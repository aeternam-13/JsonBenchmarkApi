
## API Setup (Go)
To create the go binary 

```bash
go build -o benchmark-api

./benchmark-api
```

The API can be locally tested with 
```bash
curl -s http://localhost:1313/optimal
curl -s http://localhost:1313/slower
curl -s http://localhost:1313/optimal | wc -c
curl -s http://localhost:1313/slower | wc -c
curl -v http://localhost:1313/slower 2>&1 | grep X-Duration
```

Clear the csv file (Removes all data but not headers)

```bash
curl -s http://localhost:1313/clearlog
```

Download the csv file

```bash
curl -s -O -J http://localhost:1313/csv
```
Every request will add an entry to the csv file that is used by python to generate the graphics.

Also a test file was create to measure resource usage in the functions, to compare how much more memory/CPU is used by the server.


Running test: 

```bash
go test -bench=. -benchmem
```

Output example (Arch linux):

goos: linux
goarch: amd64
pkg: jsonbenchmark
cpu: AMD Ryzen 5 5600X 6-Core Processor
BenchmarkOptimal-12    	   8613	   143528 ns/op	  59906 B/op	      3 allocs/op
BenchmarkSlower-12     	   5782	   207860 ns/op	 240129 B/op	      7 allocs/op
PASS
ok  	jsonbenchmark	2.479s


What output means:

ns/op: Nanoseconds per operation (Time).

B/op: Bytes allocated per operation (Memory Usage). This is your memory cost.

allocs/op: How many distinct objects were created in memory.

## Visualization Setup (Python)

To generate the performance charts (`benchmark_analysis.png`) from the CSV data, a python program was created, it requires a venv

### Prerequisites
- Python 3.x
- `pip`

### Create the Virtual Environment
We use a virtual environment named `stadistics` to manage dependencies.

```bash
python -m venv stadistics
```

To activate the enviroment and get the dependencies

```bash
source stadistics/bin/activate
pip install -r requirements.txt
```
Finally run the program

```bash
python stadistics.py
```