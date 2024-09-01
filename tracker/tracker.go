package tracker

import (
	"bufio"
	"bytes"
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
		"compact":    []string{"1"},
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

	// fmt.Println(url)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(body))

	r := bufio.NewReader(bytes.NewReader(body))

	decoded, _ := bencode.DecodeBencode(r)
	// fmt.Println(decoded)
	decodedMap := decoded.(map[string]bencode.BencodeValue)

	peers, err := unmarshalPeers([]byte(decodedMap["peers"].(string)))
	if err != nil {
		return nil, err
	}

	return peers, nil
}
