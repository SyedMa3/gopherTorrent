package protocol

import (
	"fmt"
	"io"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	peerID   [20]byte
}

func New(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		peerID:   peerID,
	}
}

func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buf[curr:], []byte(h.Pstr))
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.peerID[:])

	return buf
}

func Deserialize(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	pstrLen := int(lengthBuf[0])

	if pstrLen == 0 {
		err := fmt.Errorf("pstrlen cant be 0")
		return nil, err
	}

	handshakeBuf := make([]byte, 49+pstrLen)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	pstr := string(handshakeBuf[:pstrLen])
	curr := pstrLen + 8

	curr += copy(infoHash[:], handshakeBuf[curr:curr+20])
	curr += copy(peerID[:], handshakeBuf[curr:])

	h := &Handshake{
		Pstr:     pstr,
		InfoHash: infoHash,
		peerID:   peerID,
	}

	return h, nil
}
