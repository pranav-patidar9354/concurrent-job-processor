package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/pranav-patidar9354/concurrent-job-processor/internal/handlers"
)

func JobRoutes(api *gin.RouterGroup) {

	jobs := api.Group("/jobs")

	{
		jobs.POST("", handlers.CreateJob)

		jobs.GET("", handlers.GetAllJobs)

		jobs.GET("/:id", handlers.GetJobByID)

		jobs.DELETE("/:id", handlers.CancelJob)
	}

}
