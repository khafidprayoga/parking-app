﻿# Parking App

[![Go Report Card](https://goreportcard.com/badge/github.com/khafidprayoga/parking-app)](https://goreportcard.com/report/github.com/khafidprayoga/parking-app)
[![Coverage](https://img.shields.io/badge/coverage-85%25-brightgreen)](https://github.com/khafidprayoga/parking-app)

A simple parking management application developed using Go. This application allows users to manage parking spaces, park vehicles, and track parking status.

## Installation

1. Make sure Go is installed on your system (version 1.18 or newer)
2. Clone this repository:
   ```bash
   git clone https://github.com/khafidprayoga/parking-app.git
   cd parking-app
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build the application:
   ```bash
   go build -o bin/parking-app main.go
   ```

Or you can use go  package manager with this command to install as single binary   
```go install github.com/khafidprayoga/parking-app@latest```
## Running the Server

Before using the client commands, you need to start the server first:

1. Start the server:
   ```bash
   parking-app serve
   or
   bin/parking-app serve
   ```
   The server will start listening on TCP port 8080

## Usage

This application supports the following commands:

1. Create parking lot:
   ```
   parking-app create_parking_lot <number_of_slots>
   ```

2. Park vehicle:
   ```
   parking-app park <license_plate>
   ```

3. Vehicle exit:
   ```
   parking-app leave <license_plate> <duration_hours>
   ```

4. Check parking status:
   ```
   parking-app status
   ```

5. Import commands from file:
   ```
   parking-app import example/command
   ```

## Benchmark Results
v1.go output from command `task bench`
```
goos: windows
goarch: amd64
pkg: github.com/khafidprayoga/parking-app/test
cpu: AMD Ryzen 5 PRO 4650U with Radeon Graphics
BenchmarkParkingUseCase_EnterArea-12               23485             50846 ns/op             296 B/op          7 allocs/op
BenchmarkParkingUseCase_LeaveArea-12               83367             14834 ns/op             351 B/op          9 allocs/op
BenchmarkParkingUseCase_EnterAndLeave-12           70198             17949 ns/op             594 B/op         12 allocs/op
BenchmarkParkingUseCase_Parallel-12                44853             26039 ns/op             483 B/op         11 allocs/op
PASS
ok      github.com/khafidprayoga/parking-app/test       7.164s
```

v1_btree.go output from command `task bench`
```
goos: windows
goarch: amd64
pkg: github.com/khafidprayoga/parking-app/test
cpu: AMD Ryzen 5 PRO 4650U with Radeon Graphics
BenchmarkParkingUseCase_EnterArea-12              119241             10463 ns/op             299 B/op          8 allocs/op
BenchmarkParkingUseCase_LeaveArea-12              120736             10253 ns/op             342 B/op          9 allocs/op
BenchmarkParkingUseCase_EnterAndLeave-12           52960             24456 ns/op             562 B/op         12 allocs/op
BenchmarkParkingUseCase_Parallel-12                40466             28846 ns/op             489 B/op         11 allocs/op
PASS
ok      github.com/khafidprayoga/parking-app/test       6.894s

```
