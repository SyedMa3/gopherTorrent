package protocol

import (
	"encoding/binary"
	"io"
)

type messageType uint8

const (
	MsgChoke messageType = iota
	MsgUnchoke
	MsgInterested
	MsgNotInterested
	MsgHave
	MsgBitfield
	MsgRequest
	MsgPiece
	MsgCancel
)

type BitField []byte

func (bf BitField) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8

	if byteIndex < 0 || byteIndex >= len(bf) {
		return false
	}

	return (bf[byteIndex])>>(uint8(7-offset))&1 != 0
}

type Message struct {
	Type    messageType
	Payload []byte
}

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}

	length := uint32(len(m.Payload) + 1)
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.Type)
	copy(buf[5:], m.Payload)

	return buf
}

func DeserializeMessage(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, nil
	}

	msgType := messageType(messageBuf[0])
	payload := messageBuf[1:]

	m := &Message{
		Type:    msgType,
		Payload: payload,
	}

	return m, nil
}
