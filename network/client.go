package network

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/SyedMa3/gopherTorrent/bencode"
	"github.com/SyedMa3/gopherTorrent/tracker"
)

type Task struct {
	Index     int
	PieceHash [20]byte
	length    int
}

func Download(ti bencode.TorrentInfo, outPath string) {

	numPieces := len(ti.PieceHashes)
	tasksQueue := make(chan *Task, numPieces)

	for i := 0; i < numPieces; i++ {
		tasksQueue <- &Task{
			Index:     i,
			PieceHash: ti.PieceHashes[i],
			length:    ti.CalculatePieceSize(i),
		}
	}
	peers, err := tracker.GetPeers(ti)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	results := make(chan *pieceProgress, numPieces)

	AssignWork(tasksQueue, peers, ti.InfoHash, &wg, results)

	file, err := os.Create(outPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for i := 0; i < numPieces; i++ {
		progress := <-results
		_, err := file.WriteAt(progress.buf, int64(progress.index)*int64(ti.PieceLength))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Downloaded %v\n", progress.index)
		fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine()-1)
	}

	wg.Wait()
}
