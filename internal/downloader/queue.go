package downloader

import (
	"context"
	"surge/internal/messages"
	"surge/internal/utils"
	"sync"
	"time"
)

type DownloadStatus int

const (
	StatusQueued DownloadStatus = iota
	StatusDownloading
	StatusPaused // TODO: I have no idea how we are gonna manage this
	StatusCompleted
	StatusFailed

	maxDownloads = 3 // TODO: Have a max active number of downloads?
)

type QueueItem struct {
	ID      int
	Status  DownloadStatus
	Config  DownloadConfig
	Context context.Context
}

type DownloadQueue struct {
	Items        map[int]*QueueItem
	QueueList    []int
	mu           sync.Mutex
	ProgressChan chan interface{}
}

func NewDownloadQueue() *DownloadQueue {
	return &DownloadQueue{
		Items:     make(map[int]*QueueItem),
		QueueList: make([]int, 0),
	}
}

// Creates a new QueueItem and its to DownloadQueue
func (q *DownloadQueue) Add(cfg DownloadConfig) *QueueItem {

	q.mu.Lock()
	defer q.mu.Unlock()

	item := &QueueItem{
		ID:      cfg.ID,
		Config:  cfg,
		Status:  StatusQueued,
		Context: context.Background(),
	}

	q.Items[cfg.ID] = item
	q.QueueList = append(q.QueueList, item.ID)

	return item
}

func (q *DownloadQueue) ProcessQueue() {

	q.mu.Lock()
	defer q.mu.Unlock()

	active := 0
	for _, item := range q.Items {
		if item.Status == StatusDownloading {
			active++
		}
	}

	utils.Debug("Active Count: %v", active)

	for active < maxDownloads {
		nextID := -1
		for _, item := range q.Items {
			if item.Status == StatusQueued {
				nextID = item.ID
				break
			}
		}

		if nextID == -1 {
			break
		}

		item := q.Items[nextID]
		item.Status = StatusDownloading
		active++

		utils.Debug("Working on %v", nextID)

		go q.startDownload(item)
	}
}

func (q *DownloadQueue) startDownload(item *QueueItem) {

	err := TUIDownload(item.Context, item.Config)

	q.mu.Lock()
	defer q.mu.Unlock()

	if err != nil {

		item.Status = StatusFailed
		if item.Config.State != nil {
			item.Config.State.SetError(err)
		}

		if item.Config.ProgressCh != nil {
			item.Config.ProgressCh <- messages.DownloadErrorMsg{DownloadID: item.ID, Err: err}
		}
	} else {

		item.Status = StatusCompleted

		if item.Config.State != nil {
			item.Config.State.Done.Store(true)
		}

		if item.Config.ProgressCh != nil {
			item.Config.ProgressCh <- messages.DownloadCompleteMsg{DownloadID: item.ID, Elapsed: time.Since(item.Config.State.StartTime), Total: item.Config.State.TotalSize}
		}
	}

	// this avoids deadlock!
	go q.ProcessQueue()
}

func (q *DownloadQueue) getStatus(id int) DownloadStatus {

	q.mu.Lock()
	defer q.mu.Unlock()

	if item, ok := q.Items[id]; ok {
		return item.Status
	}

	return StatusFailed
}
