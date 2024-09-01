package network

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net"
	"sync"
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

func (w *Worker) sendUnchoke() error {
	msg := protocol.Message{Type: protocol.MsgUnchoke}
	_, err := w.Conn.Write(msg.Serialize())

	return err
}

func (w *Worker) sendInterested() error {
	msg := protocol.Message{Type: protocol.MsgInterested}
	_, err := w.Conn.Write(msg.Serialize())

	return err
}

func (w *Worker) sendRequest(index, begin, length int) error {
	msg := protocol.FormatRequest(index, begin, length)
	_, err := w.Conn.Write(msg.Serialize())

	return err
}

func (w *Worker) sendHave(index int) error {
	msg := protocol.FormatHave(index)
	_, err := w.Conn.Write(msg.Serialize())

	return err
}

func NewWorker(peer tracker.Peer, infoHash, peerID [20]byte) (*Worker, error) {
	// conn, err := net.DialTimeout("tcp", peer.String(), 10*time.Second)
	conn, err := net.DialTimeout("tcp", peer.String(), 15*time.Second)

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

	serialized := h.Serialize()
	_, err := conn.Write(serialized)
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

func attemptDownload(task *Task, worker Worker) (*pieceProgress, error) {
	// if worker.IsChoked {
	// 	return fmt.Errorf("peer choked")
	// }

	progress := pieceProgress{
		index:  task.Index,
		buf:    make([]byte, task.length),
		worker: &worker,
	}

	for progress.downloaded < task.length {
		if !worker.IsChoked {
			for progress.backlog < 5 && progress.requested < task.length {
				blockSize := 16384
				blockSize = min(blockSize, task.length-progress.requested)

				err := worker.sendRequest(
					progress.index,
					progress.requested,
					blockSize)
				if err != nil {
					return nil, err
				}
				progress.backlog++
				progress.requested += blockSize
			}
		}

		err := progress.readMessage()
		if err != nil {
			return nil, err
		}
	}

	return &progress, nil
}

func doWork(tasksQueue chan *Task, peer tracker.Peer, infoHash [20]byte, wg *sync.WaitGroup, results chan *pieceProgress) {
	defer wg.Done()
	worker, err := NewWorker(peer, infoHash, tracker.PeerID)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	worker.sendUnchoke()
	worker.sendInterested()

	for task := range tasksQueue {
		if !worker.BitField.HasPiece(task.Index) {
			tasksQueue <- task
			continue
		}
		//TODO
		progress, err := attemptDownload(task, *worker)
		if err != nil {
			tasksQueue <- task
			fmt.Printf("Failed to download piece %d from peer %s\n", task.Index, peer.IP)
			return
		}

		err = checkIntegrity(progress.buf, task.PieceHash)
		if err != nil {
			tasksQueue <- task
			fmt.Printf("Failed to download piece %d from peer %s\n", task.Index, peer.IP)
			continue
		}

		fmt.Printf("Successfully downloaded piece %d from peer %s with length %d\n", task.Index, peer.IP, len(progress.buf))

		worker.sendHave(task.Index)
		results <- progress

		// fmt.Printf("Downloaded %v", n)
	}
}

func AssignWork(tasksQueue chan *Task, peers []tracker.Peer, infoHash [20]byte, wg *sync.WaitGroup, results chan *pieceProgress) {
	for _, peer := range peers {
		wg.Add(1)
		go doWork(tasksQueue, peer, infoHash, wg, results)
	}
}

func checkIntegrity(piece []byte, hash [20]byte) error {
	h := sha1.New()
	h.Write(piece)
	computedHash := h.Sum(nil)
	if !bytes.Equal(computedHash, hash[:]) {
		return fmt.Errorf("expected %v, got: %v", hash, computedHash)
	}

	return nil
}
