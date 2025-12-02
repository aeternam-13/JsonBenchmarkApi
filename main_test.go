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

/*

Running tool: /opt/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkSlower$ jsonbenchmark

goos: linux
goarch: amd64
pkg: jsonbenchmark
cpu: AMD Ryzen 5 5600X 6-Core Processor
BenchmarkSlower-12    	    5613	    208369 ns/op	  240129 B/op	       7 allocs/op
PASS
ok  	jsonbenchmark	1.196s

Running tool: /opt/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkOptimal$ jsonbenchmark

goos: linux
goarch: amd64
pkg: jsonbenchmark
cpu: AMD Ryzen 5 5600X 6-Core Processor
BenchmarkOptimal-12    	    8821	    144617 ns/op	   59906 B/op	       3 allocs/op
PASS
ok  	jsonbenchmark	1.294s

go test -bench=. -benchmem


goos: linux
goarch: amd64
pkg: jsonbenchmark
cpu: AMD Ryzen 5 5600X 6-Core Processor
BenchmarkOptimal-12    	   8613	   143528 ns/op	  59906 B/op	      3 allocs/op
BenchmarkSlower-12     	   5782	   207860 ns/op	 240129 B/op	      7 allocs/op
PASS
ok  	jsonbenchmark	2.479s


ns/op: Nanoseconds per operation (Time).

B/op: Bytes allocated per operation (Memory Usage). This is your memory cost.

allocs/op: How many distinct objects were created in memory.
*/
