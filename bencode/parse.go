package bencode

import "crypto/sha1"

type bencodeInfo struct {
	pieces      string
	pieceLength int64
	length      int64
	name        string
}

type bencodeTorrent struct {
	announce string
	info     bencodeInfo
}

type TorrentInfo struct {
	announce    string
	infoHash    [20]byte
	pieceHashes [][20]byte
	pieceLength int64
	length      int64
	name        string
}

type BencodeValue interface{}

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

func (bto bencodeTorrent) toTorrentInfo() TorrentInfo {
	h := sha1.New()

	ti := TorrentInfo{
		announce:    bto.announce,
		infoHash:    [20]byte(h.Sum(encodeBencode(bto.info))),
		pieceHashes: toHashes(bto.info.pieces),
		pieceLength: bto.info.pieceLength,
		length:      bto.info.length,
		name:        bto.info.name,
	}

	return ti
}
