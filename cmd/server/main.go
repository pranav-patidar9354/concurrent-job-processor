package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/pranav-patidar9354/concurrent-job-processor/config"
	"github.com/pranav-patidar9354/concurrent-job-processor/internal/docs"
	"github.com/pranav-patidar9354/concurrent-job-processor/internal/routes"
	"github.com/pranav-patidar9354/concurrent-job-processor/internal/worker"
)

func main() {

	// Load environment variables
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

	// Connect to database
	config.ConnectDatabase()

	// Start worker pool
	worker.StartWorkerPool(3)

	// Create Gin router
	router := gin.Default()

	// Swagger Documentation (FastAPI style /docs)
	docs.RegisterDocsRoutes(router)

	// Health check
	router.GET("/", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Concurrent Job Processor API is running",
		})
	})

	// API routes
	api := router.Group("/api/v1")

	routes.JobRoutes(api)

	// Get server port (PORT for Render/Cloud, SERVER_PORT for local)
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("SERVER_PORT")
	}
	if port == "" {
		port = "8081"
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start HTTP server in a goroutine
	go func() {

		fmt.Println(
			"Server running on port:",
			port,
		)

		err := server.ListenAndServe()

		if err != nil &&
			err != http.ErrServerClosed {

			panic(
				"Failed to start server",
			)
		}
	}()

	// Wait for shutdown signal
	quit := make(
		chan os.Signal,
		1,
	)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	fmt.Println(
		"Shutdown signal received",
	)

	// Give HTTP server maximum 10 seconds to shut down
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	// Stop accepting new HTTP requests
	err = server.Shutdown(ctx)

	if err != nil {

		fmt.Println(
			"HTTP server forced to shutdown:",
			err,
		)
	}

	fmt.Println(
		"HTTP server stopped",
	)

	// Close job queue and wait for workers
	worker.StopWorkerPool()

	fmt.Println(
		"Application shutdown completed",
	)
}