package worker

import (
	"context"
	"sync"
)

var (
	jobContexts = make(map[uint]context.CancelFunc)
	contextMu   sync.Mutex
)

func RegisterJobContext(jobID uint, cancel context.CancelFunc) {

	contextMu.Lock()
	defer contextMu.Unlock()

	jobContexts[jobID] = cancel
}

func CancelJobContext(jobID uint) bool {

	contextMu.Lock()
	defer contextMu.Unlock()

	cancel, exists := jobContexts[jobID]

	if !exists {
		return false
	}

	cancel()

	return true
}

func RemoveJobContext(jobID uint) {

	contextMu.Lock()
	defer contextMu.Unlock()

	delete(jobContexts, jobID)
}