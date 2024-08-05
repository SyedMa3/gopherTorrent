package connection

import (
	"github.com/SyedMa3/gopherTorrent/bencode"
	"github.com/SyedMa3/gopherTorrent/tracker"
)

type Client struct {
	PeerID string
	Peers  []tracker.Peer
	Ti     bencode.TorrentInfo
}
