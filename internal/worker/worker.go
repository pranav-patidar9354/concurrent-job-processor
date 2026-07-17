package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/pranav-patidar9354/concurrent-job-processor/config"
	"github.com/pranav-patidar9354/concurrent-job-processor/internal/models"
)

func StartWorker(workerID int) {

	fmt.Printf("Worker %d started\n", workerID)

	for job := range JobQueue {

		var currentJob models.Job

		err := config.DB.First(
			&currentJob,
			job.ID,
		).Error

		if err != nil {
			continue
		}

		if currentJob.Status == "cancelled" {

			fmt.Printf(
				"Worker %d skipped cancelled Job %d\n",
				workerID,
				job.ID,
			)

			continue
		}

		fmt.Printf(
			"Worker %d picked Job %d (%s)\n",
			workerID,
			job.ID,
			job.Type,
		)

		processJob(workerID, job)
	}

	fmt.Printf("Worker %d stopped\n", workerID)
}

func processJob(workerID int, job models.Job) {

	// Each job gets a maximum processing time of 10 seconds
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	RegisterJobContext(job.ID, cancel)

	defer func() {
		cancel()
		RemoveJobContext(job.ID)
	}()

	startedAt := time.Now()

	err := config.DB.Model(&job).Updates(map[string]interface{}{
		"status":     "processing",
		"started_at": &startedAt,
	}).Error

	if err != nil {
		fmt.Printf(
			"Worker %d failed to update Job %d status\n",
			workerID,
			job.ID,
		)
		return
	}

	fmt.Printf(
		"Worker %d processing Job %d (%s)\n",
		workerID,
		job.ID,
		job.Type,
	)

	var duration time.Duration
	var result string

	switch job.Type {

	case "generate_report":

		duration = 6 * time.Second

		result = fmt.Sprintf(
			"Report generated successfully: %s",
			job.Payload,
		)

	case "send_email":

		duration = 3 * time.Second

		result = fmt.Sprintf(
			"Email sent successfully: %s",
			job.Payload,
		)

	case "process_data":

		duration = 8 * time.Second

		result = fmt.Sprintf(
			"Data processed successfully: %s",
			job.Payload,
		)

	default:

		config.DB.Model(&job).Updates(map[string]interface{}{
			"status":        "failed",
			"error_message": fmt.Sprintf(
				"unsupported job type: %s",
				job.Type,
			),
		})

		fmt.Printf(
			"Worker %d failed Job %d: unsupported job type\n",
			workerID,
			job.ID,
		)

		return
	}

	select {

	case <-time.After(duration):

		// Job finished normally

	case <-ctx.Done():

		status := "cancelled"
		errorMessage := "job was cancelled"

		if ctx.Err() == context.DeadlineExceeded {
			status = "failed"
			errorMessage = "job processing timed out"
		}

		config.DB.Model(&job).Updates(map[string]interface{}{
			"status":        status,
			"error_message": errorMessage,
		})

		fmt.Printf(
			"Worker %d stopped Job %d: %s\n",
			workerID,
			job.ID,
			errorMessage,
		)

		return
	}

	completedAt := time.Now()

	err = config.DB.Model(&job).Updates(map[string]interface{}{
		"status":        "completed",
		"result":        result,
		"completed_at":  &completedAt,
		"error_message": "",
	}).Error

	if err != nil {

		fmt.Printf(
			"Worker %d failed to save result for Job %d\n",
			workerID,
			job.ID,
		)

		return
	}

	fmt.Printf(
		"Worker %d completed Job %d (%s)\n",
		workerID,
		job.ID,
		job.Type,
	)
}