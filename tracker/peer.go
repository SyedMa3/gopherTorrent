package tracker

import (
	"encoding/binary"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
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

func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}
