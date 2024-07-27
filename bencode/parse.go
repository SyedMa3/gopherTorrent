package bencode

import (
	"crypto/sha1"
)

type bencodeInfo struct {
	length      int64
	name        string
	pieceLength int64
	pieces      string
}

type BencodeTorrent struct {
	announce string
	info     bencodeInfo
}

type TorrentInfo struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int64
	Length      int64
	Name        string
}

type BencodeValue interface{}

func NewBencodeTorrent(result map[string]BencodeValue) *BencodeTorrent {
	bto := &BencodeTorrent{
		announce: result["announce"].(string),
		info: bencodeInfo{
			pieces:      result["info"].(map[string]BencodeValue)["pieces"].(string),
			pieceLength: result["info"].(map[string]BencodeValue)["piece length"].(int64),
			length:      result["info"].(map[string]BencodeValue)["length"].(int64),
			name:        result["info"].(map[string]BencodeValue)["name"].(string),
		},
	}

	return bto
}

func toHashes(pieces string) [][20]byte {
	buf := []byte(pieces)

	num := len(buf) / 20
	hashes := make([][20]byte, num+1)

	for i := 0; i < num; i++ {
		copy(hashes[i][:], buf[i*20:(i+1)*20])
	}
	if len(buf)%20 != 0 {
		copy(hashes[num][:], buf[(num+1)*20:])
	}

	return hashes
}

func (bto BencodeTorrent) ToTorrentInfo() TorrentInfo {
	encodedInfo := encodeInfo(bto.info)

	// fmt.Println("ss", (sha1.Sum(encodedInfo)))

	ti := TorrentInfo{
		Announce:    bto.announce,
		InfoHash:    [20]byte(sha1.Sum(encodedInfo)),
		PieceHashes: toHashes(bto.info.pieces),
		PieceLength: bto.info.pieceLength,
		Length:      bto.info.length,
		Name:        bto.info.name,
	}

	// fmt.Println(string(h.Sum(encodeBencode(bto.info))))

	return ti
}
