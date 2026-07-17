package services

import (
	"fmt"

	"github.com/pranav-patidar9354/concurrent-job-processor/internal/models"
	"github.com/pranav-patidar9354/concurrent-job-processor/internal/repository"
	"github.com/pranav-patidar9354/concurrent-job-processor/internal/worker"
)

func CreateJob(req models.CreateJobRequest) (*models.Job, error) {

	job := models.Job{
		Type:    req.Type,
		Payload: req.Payload,
		Status:  "queued",
	}

	err := repository.CreateJob(&job)

	if err != nil {
		return nil, err
	}

	worker.AddJob(job)

	return &job, nil
}

func GetJobByID(id uint) (*models.Job, error) {
	return repository.GetJobByID(id)
}

func GetAllJobs() ([]models.Job, error) {
	return repository.GetAllJobs()
}

func CancelJob(id uint) error {

	job, err := repository.GetJobByID(id)

	if err != nil {
		return err
	}

	if job.Status == "completed" ||
		job.Status == "failed" ||
		job.Status == "cancelled" {

		return fmt.Errorf(
			"job cannot be cancelled because its status is %s",
			job.Status,
		)
	}

	cancelled := worker.CancelJobContext(id)

	if !cancelled {

		// Job may still be waiting inside the queue.
		err = repository.UpdateJobStatus(
			id,
			"cancelled",
		)

		return err
	}

	return nil
}