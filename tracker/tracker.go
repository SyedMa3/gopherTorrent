package tracker

import (
	"crypto/rand"
	"io"
	"net/http"
	"net/url"
	"strconv"

	bencode "github.com/SyedMa3/gopherTorrent/bencode"
)

var PeerID [20]byte

const Port uint16 = 6881

func init() {
	rand.Read(PeerID[:])
}

func buildTrackerURL(t bencode.TorrentInfo) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(PeerID[:])},
		"port":       []string{strconv.Itoa(int(Port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(int(t.Length))},
	}

	base.RawQuery = params.Encode()

	return base.String(), nil
}

// Returns the list of Peers given the TorrentInfo
func GetPeers(t bencode.TorrentInfo) ([]Peer, error) {
	url, err := buildTrackerURL(t)
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
