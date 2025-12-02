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

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_:;{}¿?="

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
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

	// Initialize Logger (function is now in logger.go)
	logger, err := NewBenchmarkLogger("benchmark_data.csv")
	if err != nil {
		panic(err)
	}
	fmt.Println("Logging data to benchmark_data.csv...")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(NetworkMonitorMiddleware(logger))

	// The target size for the String that will be compared in the endpoints and android app

	r.GET("/optimal", func(c *gin.Context) {
		response := OptimalParsing()
		c.JSON(http.StatusOK, response)
	})

	r.GET("/slower", func(c *gin.Context) {
		response := SlowerParsing()

		c.JSON(http.StatusOK, response)
	})

	r.Run(":1313")
}
