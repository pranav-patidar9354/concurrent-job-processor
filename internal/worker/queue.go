package worker

import (
	"github.com/pranav-patidar9354/concurrent-job-processor/internal/models"
)

var JobQueue = make(chan models.Job, 100)

func AddJob(job models.Job) {

	JobQueue <- job
}