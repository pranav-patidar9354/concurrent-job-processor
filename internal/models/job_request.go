package models

type CreateJobRequest struct {
	Type    string `json:"type" binding:"required"`
	Payload string `json:"payload" binding:"required"`
}