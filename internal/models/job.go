package models

import "time"

type Job struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Type string `json:"type" gorm:"not null"`

	Payload string `json:"payload" gorm:"type:text"`

	Status string `json:"status" gorm:"not null;default:queued"`

	Result string `json:"result" gorm:"type:text"`

	ErrorMessage string `json:"error_message" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	StartedAt *time.Time `json:"started_at"`

	CompletedAt *time.Time `json:"completed_at"`
}