package network

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/SyedMa3/gopherTorrent/protocol"
	"github.com/SyedMa3/gopherTorrent/tracker"
)

type Worker struct {
	Conn     net.Conn
	IsChoked bool
	peer     tracker.Peer
	infoHash [20]byte
	peerID   [20]byte
	BitField protocol.BitField
}

func NewWorker(peer tracker.Peer, infoHash, peerID [20]byte) (*Worker, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	err = doHandshake(conn, infoHash, peerID)
	if err != nil {
		return nil, err
	}

	bitFieldMsg, err := protocol.DeserializeMessage(conn)
	if err != nil {
		return nil, err
	}

	if bitFieldMsg.Type != protocol.MsgBitfield {
		return nil, fmt.Errorf("expected %v, got: %v", protocol.MsgBitfield, bitFieldMsg.Type)
	}

	return &Worker{
		Conn:     conn,
		IsChoked: true,
		peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
		BitField: bitFieldMsg.Payload,
	}, nil
}

func doHandshake(conn net.Conn, infoHash, peerID [20]byte) error {
	h := protocol.NewHandshake(infoHash, peerID)

	_, err := conn.Write(h.Serialize())
	if err != nil {
		return err
	}

	resp, err := protocol.DeserializeHandshake(conn)
	if err != nil {
		return err
	}

	if !bytes.Equal(resp.InfoHash[:], infoHash[:]) {
		return fmt.Errorf("expected %v, Got: %v", infoHash, resp.InfoHash)
	}

	return nil
}
