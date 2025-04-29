package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func main() {
	router := gin.Default()

	router.GET("/temperature", func(c *gin.Context) {
		location := c.DefaultQuery("location", "default")

		if location == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"location":    location,
			"value":       -50 + rand.Float64()*100,
			"unit":        "unit",
			"status":      "ok",
			"timestamp":   time.Now(),
			"description": "super temperature",
		})
	})

	srv := &http.Server{
		Addr:    getEnv("ENDPOINT", ":8081"),
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited properly")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
