package repository

import (
	"github.com/pranav-patidar9354/concurrent-job-processor/config"
	"github.com/pranav-patidar9354/concurrent-job-processor/internal/models"
)

func CreateJob(job *models.Job) error {
	return config.DB.Create(job).Error
}

func GetJobByID(id uint) (*models.Job, error) {

	var job models.Job

	err := config.DB.First(&job, id).Error

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func GetAllJobs() ([]models.Job, error) {

	var jobs []models.Job

	err := config.DB.Order("created_at DESC").Find(&jobs).Error

	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func UpdateJobStatus(id uint, status string) error {

	return config.DB.Model(&models.Job{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}