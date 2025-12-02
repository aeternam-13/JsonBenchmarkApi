package main

import (
	"encoding/base64"
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
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	r := gin.Default()

	// Clair Obscure GOTY
	commonFoo := "For those who come after !!!"

	// The target size for the String that will be compared in the endpoints and android app
	const targetSize = 2000

	r.GET("/optimal", func(c *gin.Context) {
		// Generate the raw large string
		largeString := generateRandomString(targetSize)

		response := DummyClass{
			Foo1:   commonFoo,
			Foo2:   commonFoo,
			Foo3:   commonFoo,
			Foo4:   commonFoo,
			Foo5:   commonFoo,
			Foo6:   commonFoo,
			Target: largeString, // Raw 2000 chars
		}

		c.JSON(http.StatusOK, response)
	})

	r.GET("/slower", func(c *gin.Context) {
		// Generate the raw large string
		originalString := generateRandomString(targetSize)

		// First encode
		encodedOnce := base64.StdEncoding.EncodeToString([]byte(originalString))

		// Second encode
		encodedTwice := base64.StdEncoding.EncodeToString([]byte(encodedOnce))

		response := DummyClass{
			Foo1:   commonFoo,
			Foo2:   commonFoo,
			Foo3:   commonFoo,
			Foo4:   commonFoo,
			Foo5:   commonFoo,
			Foo6:   commonFoo,
			Target: encodedTwice, // Encoded twice
		}

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

*/
