package tracker

import (
	"encoding/binary"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"

	bencode "github.com/SyedMa3/gopherTorrent/bencode"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

// func main() {
// 	filename := "debian-12.6.0-amd64-netinst.iso.torrent"
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		fmt.Printf("Error opening file: %v\n", err)
// 		return
// 	}
// 	defer file.Close()

// 	reader := bufio.NewReader(file)
// 	result, err := bencode.DecodeBencode(reader)
// 	if err != nil {
// 		fmt.Printf("Error decoding bencode: %v\n", err)
// 		return
// 	}

// 	newResult := result.(map[string]bencode.BencodeValue)

// 	// fmt.Println(reflect.TypeOf(newResult["info"]))
// 	// fmt.Println(newResult["info"])

// 	bto := bencode.NewBencodeTorrent(newResult)

// 	// fmt.Println(bto)

// 	ti := bto.ToTorrentInfo()

// 	peerID := "js8uJhsyw64mKJi9tyRa"

// 	port := uint16(6881)

// 	peers, _ := getPeers(ti, [20]byte([]byte(peerID)), port)
// 	fmt.Println(peers)
// }

func buildTrackerURL(t bencode.TorrentInfo, peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(int(t.Length))},
	}

	base.RawQuery = params.Encode()

	return base.String(), nil
}

// Returns the list of Peers given the TorrentInfo
func GetPeers(t bencode.TorrentInfo, peerID [20]byte, port uint16) ([]Peer, error) {
	url, err := buildTrackerURL(t, peerID, port)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	peers, err := unmarshalPeers(body)
	if err != nil {
		return nil, err
	}

	return peers, nil
}

func unmarshalPeers(peers []byte) ([]Peer, error) {
	peerSize := 6
	numPeers := len(peers) / peerSize

	unmarshalledPeers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		offset := peerSize * i
		unmarshalledPeers[i] = Peer{
			IP:   net.IP(peers[offset : offset+4]),
			Port: binary.BigEndian.Uint16(peers[offset+4 : offset+6]),
		}
	}

	return unmarshalledPeers, nil
}
