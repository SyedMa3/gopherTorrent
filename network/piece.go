package network

import (
	"github.com/SyedMa3/gopherTorrent/protocol"
)

type pieceProgress struct {
	index      int
	buf        []byte
	downloaded int
	backlog    int
	requested  int
	worker     *Worker
}

func (p *pieceProgress) readMessage() error {
	msg, err := protocol.DeserializeMessage(p.worker.Conn)
	if err != nil {
		return err
	}

	if msg == nil {
		return nil
	}

	switch msg.Type {
	case protocol.MsgChoke:
		p.worker.IsChoked = true
	case protocol.MsgUnchoke:
		p.worker.IsChoked = false
	case protocol.MsgHave:
		index, err := protocol.ParseHave(msg)
		if err != nil {
			return err
		}
		p.worker.BitField.SetPiece(index)
	case protocol.MsgPiece:
		// fmt.Printf("received piece %d\n", p.index)
		n, err := protocol.ParsePiece(p.index, p.buf, msg)
		if err != nil {
			return err
		}
		p.downloaded += n
		p.backlog--
	}
	return nil
}
