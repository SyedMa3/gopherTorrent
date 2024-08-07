package network

import (
	"github.com/SyedMa3/gopherTorrent/bencode"
)

type Task struct {
	Index     int
	PieceHash [20]byte
}

func Download(ti bencode.TorrentInfo) {

	numPieces := len(ti.PieceHashes)

	tasksQueue := make(chan *Task, numPieces)

	for i := 0; i < numPieces; i++ {
		tasksQueue <- &Task{
			Index:     i,
			PieceHash: ti.PieceHashes[i],
		}
	}
	// peers, err := tracker.GetPeers(ti)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}

//TODO: start download for every worker
