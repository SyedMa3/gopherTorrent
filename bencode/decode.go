package bencode

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// func main() {
// 	filename := "debian-12.6.0-amd64-netinst.iso.torrent"
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		fmt.Printf("Error opening file: %v\n", err)
// 		return
// 	}
// 	defer file.Close()

// 	reader := bufio.NewReader(file)
// 	result, err := decodeBencode(reader)
// 	if err != nil {
// 		fmt.Printf("Error decoding bencode: %v\n", err)
// 		return
// 	}

// 	newResult := result.(map[string]BencodeValue)

// 	bto := &bencodeTorrent{
// 		announce: newResult["announce"].(string),
// 		info: bencodeInfo{
// 			pieces:      newResult["info"].(map[string]BencodeValue)["pieces"].(string),
// 			pieceLength: newResult["info"].(map[string]BencodeValue)["piece length"].(int64),
// 			length:      newResult["info"].(map[string]BencodeValue)["length"].(int64),
// 			name:        newResult["info"].(map[string]BencodeValue)["name"].(string),
// 		},
// 	}

// 	// fmt.Println(bto)

// 	ti, _ := bto.toTorrentInfo()

// 	fmt.Printf("Parsed torrentfile: %+v\n", ti)
// }

func decodeBencode(r *bufio.Reader) (BencodeValue, error) {
	b, err := r.Peek(1)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(b[0]))

	switch b[0] {
	case 'i':
		return decodeInt(r)
	case 'l':
		return decodeList(r)
	case 'd':
		return decodeDict(r)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return decodeString(r)
	default:
		return nil, fmt.Errorf("unknown bencode type: %c", b[0])
	}
}

func decodeInt(r *bufio.Reader) (int64, error) {
	r.ReadByte() // consume 'i'
	var num string
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		if b == 'e' {
			break
		}
		num += string(b)
	}
	return strconv.ParseInt(num, 10, 64)
}

func decodeString(r *bufio.Reader) (string, error) {
	var length string
	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if b == ':' {
			break
		}
		length += string(b)
	}
	l, err := strconv.ParseInt(length, 10, 64)
	// fmt.Println(l)
	if err != nil {
		return "", err
	}
	bytes := make([]byte, l)
	_, err = io.ReadFull(r, bytes)
	return string(bytes), err
}

func decodeList(r *bufio.Reader) ([]BencodeValue, error) {
	r.ReadByte() // consume 'l'
	var list []BencodeValue
	for {
		b, err := r.Peek(1)
		if err != nil {
			return nil, err
		}
		if b[0] == 'e' {
			r.ReadByte() // consume 'e'
			break
		}
		value, err := decodeBencode(r)
		if err != nil {
			return nil, err
		}
		list = append(list, value)
	}
	return list, nil
}

func decodeDict(r *bufio.Reader) (map[string]BencodeValue, error) {
	r.ReadByte() // consume 'd'
	dict := make(map[string]BencodeValue)
	for {
		b, err := r.Peek(1)
		// fmt.Println(string(b))
		if err != nil {
			return nil, err
		}
		if b[0] == 'e' {
			r.ReadByte() // consume 'e'
			break
		}
		key, err := decodeString(r)
		if err != nil {
			return nil, err
		}
		value, err := decodeBencode(r)
		if err != nil {
			return nil, err
		}
		dict[key] = value
	}
	return dict, nil
}
