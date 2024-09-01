package main

import (
	"fmt"

	"github.com/SyedMa3/gopherTorrent/bencode"
	"github.com/SyedMa3/gopherTorrent/network"
)

func main() {
	// args := os.Args[1:]

	// if len(args) != 1 {
	// 	fmt.Printf("Incorrent number of arguemnts\n")
	// 	return
	// }

	// filePath := args[0]
	filePath := "./debian-iso.torrent"
	outPath := "./debian-iso.iso"
	ti, err := bencode.FileToTorrentInfo(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	network.Download(*ti, outPath)

}
