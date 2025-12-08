package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Clair Obscure GOTY
const commonFoo = "For those who come after !!!"
const targetSize = 25000
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_:;{}?="

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func OptimalParsing() DummyClass {
	largeString := generateRandomString(targetSize)
	return DummyClass{
		Foo1: commonFoo, Foo2: commonFoo, Foo3: commonFoo,
		Foo4: commonFoo, Foo5: commonFoo, Foo6: commonFoo,
		Target: InnerTarget{
			Foo1: commonFoo,
			Foo2: commonFoo,
			Foo3: commonFoo,
			Target: DeeperTarget{
				Foo1:         commonFoo,
				Foo2:         commonFoo,
				Foo3:         commonFoo,
				TargetString: largeString,
			},
		},
	}
}

func SlowerParsing() DummySlowerClass {
	targetResult := ""
	originalString := generateRandomString(targetSize)

	deeperTarget := DeeperTarget{
		Foo1:         commonFoo,
		Foo2:         commonFoo,
		Foo3:         commonFoo,
		TargetString: originalString,
	}

	deeperTargetEncoded, err := json.Marshal(deeperTarget)

	if err != nil {
		fmt.Println("Error:", err)
		targetResult = ""
	}

	encodedOnce := base64.StdEncoding.EncodeToString(deeperTargetEncoded)

	innerTarget := InnerSlowerTarget{
		Foo1:   commonFoo,
		Foo2:   commonFoo,
		Foo3:   commonFoo,
		Target: encodedOnce,
	}

	innerTargetEncoded, err := json.Marshal(innerTarget)

	if err != nil {
		fmt.Println("Error:", err)
		targetResult = ""
	}

	targetResult = base64.StdEncoding.EncodeToString([]byte(innerTargetEncoded))

	return DummySlowerClass{
		Foo1: commonFoo, Foo2: commonFoo, Foo3: commonFoo,
		Foo4: commonFoo, Foo5: commonFoo, Foo6: commonFoo,
		Target: targetResult,
	}
}
