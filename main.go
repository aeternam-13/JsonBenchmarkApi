package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize Logger (function is now in customLogger.go)
	customLogger, err := NewBenchmarkLogger("benchmark_data.csv")
	if err != nil {
		panic(err)
	}
	fmt.Println("Logging data to benchmark_data.csv...")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(NetworkMonitorMiddleware(customLogger))

	// The target size for the String that will be compared in the endpoints and android app

	r.GET("/optimal", func(c *gin.Context) {
		response := OptimalParsing()
		c.JSON(http.StatusOK, response)
	})

	r.GET("/slower", func(c *gin.Context) {
		response := SlowerParsing()

		c.JSON(http.StatusOK, response)
	})

	r.GET("/clearlog", func(c *gin.Context) {
		if err := customLogger.Clear(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Log cleared"})
	})

	r.GET("csv", func(c *gin.Context) {
		customLogger.DownloadHandler(c)
	})

	r.Run(":1313")
}
