package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/pranav-patidar9354/concurrent-job-processor/internal/models"
	"github.com/pranav-patidar9354/concurrent-job-processor/internal/services"
	"github.com/pranav-patidar9354/concurrent-job-processor/pkg/response"
)

func CreateJob(c *gin.Context) {

	var req models.CreateJobRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"type and payload are required",
		)
		return
	}

	job, err := services.CreateJob(req)

	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"Job submitted successfully",
		job,
	)
}

func GetJobByID(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"Invalid job ID",
		)
		return
	}

	job, err := services.GetJobByID(uint(id))

	if err != nil {

		if err == gorm.ErrRecordNotFound {
			response.Error(
				c,
				http.StatusNotFound,
				"Job not found",
			)
			return
		}

		response.InternalServerError(c)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Job fetched successfully",
		job,
	)
}

func GetAllJobs(c *gin.Context) {

	jobs, err := services.GetAllJobs()

	if err != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Jobs fetched successfully",
		jobs,
	)
}

func CancelJob(c *gin.Context) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			"Invalid job ID",
		)

		return
	}

	err = services.CancelJob(uint(id))

	if err != nil {

		if err == gorm.ErrRecordNotFound {

			response.Error(
				c,
				http.StatusNotFound,
				"Job not found",
			)

			return
		}

		response.Error(
			c,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Job cancellation requested successfully",
		nil,
	)
}