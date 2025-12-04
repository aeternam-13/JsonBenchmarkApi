package main

import (
	"testing"
)

// Benchmark for Optimal Logic
func BenchmarkOptimal(b *testing.B) {
	// Run the function b.N times
	for i := 0; i < b.N; i++ {
		OptimalParsing()
	}
}

// Benchmark for Slower Logic
func BenchmarkSlower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SlowerParsing()
	}
}
