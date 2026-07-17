package worker

import (
	"fmt"
	"sync"
)

var workerWG sync.WaitGroup

func StartWorkerPool(numberOfWorkers int) {

	fmt.Printf(
		"Starting worker pool with %d workers\n",
		numberOfWorkers,
	)

	for i := 1; i <= numberOfWorkers; i++ {

		workerWG.Add(1)

		go func(workerID int) {

			defer workerWG.Done()

			StartWorker(workerID)

		}(i)
	}
}

func StopWorkerPool() {

	fmt.Println("Stopping worker pool...")

	close(JobQueue)

	workerWG.Wait()

	fmt.Println("All workers stopped successfully")
}