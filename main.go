package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Struct to match Kotlin class
type DummyClass struct {
	Foo1   string `json:"foo1"`
	Foo2   string `json:"foo2"`
	Foo3   string `json:"foo3"`
	Foo4   string `json:"foo4"`
	Foo5   string `json:"foo5"`
	Foo6   string `json:"foo6"`
	Target string `json:"target"`
}

// Helper: Generates a random string of fixed length
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_:;{}¿?="

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Method to calculate the time spent in creating the response, works for both
func NetworkMonitorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // Process request

		duration := time.Since(start)
		responseSize := c.Writer.Size() // <--- Get response size in Bytes

		// Add metrics to headers so the client can see them
		c.Writer.Header().Set("X-Duration", fmt.Sprintf("%v", duration))
		c.Writer.Header().Set("X-Size-Bytes", fmt.Sprintf("%d", responseSize))

		// Log formatted output for your "Showcase"
		fmt.Printf("path: %-10s | time: %-12v | size: %d bytes\n",
			c.Request.URL.Path, duration, responseSize)
	}
}

// Clair Obscure GOTY
const commonFoo = "For those who come after !!!"
const targetSize = 25000

func OptimalParsing() DummyClass {
	largeString := generateRandomString(targetSize)
	return DummyClass{
		Foo1: commonFoo, Foo2: commonFoo, Foo3: commonFoo,
		Foo4: commonFoo, Foo5: commonFoo, Foo6: commonFoo,
		Target: largeString,
	}
}

func SlowerParsing() DummyClass {
	originalString := generateRandomString(targetSize)
	encodedOnce := base64.StdEncoding.EncodeToString([]byte(originalString))
	encodedTwice := base64.StdEncoding.EncodeToString([]byte(encodedOnce))

	return DummyClass{
		Foo1: commonFoo, Foo2: commonFoo, Foo3: commonFoo,
		Foo4: commonFoo, Foo5: commonFoo, Foo6: commonFoo,
		Target: encodedTwice,
	}
}

func main() {
	r := gin.Default()
	r.Use(NetworkMonitorMiddleware())

	// The target size for the String that will be compared in the endpoints and android app

	r.GET("/optimal", func(c *gin.Context) {
		response := OptimalParsing()
		c.JSON(http.StatusOK, response)
	})

	r.GET("/slower", func(c *gin.Context) {
		response := SlowerParsing()

		c.JSON(http.StatusOK, response)
	})

	r.Run(":8080")
}

/**


COMPILE ->

go build -o benchmark-api main.go

Execute ->

./benchmark-api

Test ->

curl -s http://localhost:8080/optimal | wc -c
curl -s http://localhost:8080/slower | wc -c


curl -v http://localhost:8080/slower 2>&1 | grep X-Duration

*/
