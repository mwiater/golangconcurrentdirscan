# golangconncurentdirscan

## Article
This code is reference for the following Medium article: [Extending Golang’s Standard Packages with Concurrency: Winning Big By Starting Small](https://levelup.gitconnected.com/extending-golangs-standard-packages-with-concurrency-winning-big-by-starting-small-8ba5e5ae163c)

## Summary

This Go program, `golangconncurentdirscan`, provides a robust and efficient way to scan a directory tree. It implements two methods for scanning directories: one using the `filepath.WalkDir` function and the other using concurrent goroutines. The `filepath.WalkDir` method performs a sequential directory scan, while the goroutine-based method leverages concurrent programming to scan directories in parallel, potentially reducing the overall scanning time.

The program compares the performance of these two methods by measuring execution time, memory usage, and the number of files and directories found. It then provides a detailed comparison of the results, including any discrepancies between the two methods. This comparison is useful for understanding the trade-offs between sequential and concurrent directory scanning in terms of speed and resource consumption.

## How to Run the Code

To run the code, ensure you have Go installed on your system. You can run the program directly from the source code using the `go run` command. Open your terminal or command prompt, navigate to the directory containing the `main.go` file, and execute the following command:

Linux:
```sh
go run . -path="/home/matt/projects" -concurrencyMultiplier=1
```

Windows:
```sh
go run . -path="C:\Users\Matt\projects" -concurrencyMultiplier=1
```

This command compiles and runs the Go program, performing the directory scan and displaying the results.

## How to Build the Code

### Linux

To build the code on a Linux system, use the `go build` command to compile the program into an executable. Open your terminal, navigate to the directory containing the `main.go` file, and run the following command:

```sh
go build -o bin/scan
```

This command compiles the Go program and outputs the executable to the `bin` directory with the name `scan`. You can then run the executable using:

```sh
./bin/scan
```

### Windows

To build the code on a Windows system, use the `go build` command to compile the program into an executable. Open your command prompt, navigate to the directory containing the `main.go` file, and run the following command:

```sh
go build -o bin\scan.exe
```

This command compiles the Go program and outputs the executable to the `bin` directory with the name `scan.exe`. Open your command prompt, navigate to the directory containing the `main.go` file, and run the following command:

```sh
bin\scan.exe
```

By following these instructions, you can easily compile and run the `golangconncurentdirscan` program on both Linux and Windows systems, allowing you to compare the performance of sequential and concurrent directory scanning methods.

## Example Output


 ```
  Baseline: No Concurrency | Directory Scan Comparison: /home/matt/projects
 FUNCTION         NUMCPUS  CONCURRENCY   FILES  DIRECTORIES  OBJECTS         SIZE     MEMORY  EXECUTION TIME  SPEED INCREASE
 WalkDir          N/A      N/A          477859        85515   563374  11650846047  482278744    1.754791989s  1.3x
 Concurrent Read  8        N/A          477859        85515   563374  11650846047  519758552    2.308008297s  1.0x


 Test Number: 1 | Directory Scan Comparison: /home/matt/projects
 FUNCTION         NUMCPUS  CONCURRENCY   FILES  DIRECTORIES  OBJECTS         SIZE     MEMORY  EXECUTION TIME  SPEED INCREASE
 WalkDir          N/A      N/A          477859        85515   563374  11650846047  482193672    1.895348317s  1.2x
 Concurrent Read  8        1            477859        85515   563374  11650846047  511525136    2.237825253s  1.0x


 Test Number: 2 | Directory Scan Comparison: /home/matt/projects
 FUNCTION         NUMCPUS  CONCURRENCY   FILES  DIRECTORIES  OBJECTS         SIZE     MEMORY  EXECUTION TIME  SPEED INCREASE
 WalkDir          N/A      N/A          477859        85515   563374  11650846047  482159056    1.783704648s  1.0x
 Concurrent Read  8        2            477859        85515   563374  11650846047  495747136    1.166968233s  1.5x


 Test Number: 3 | Directory Scan Comparison: /home/matt/projects
 FUNCTION         NUMCPUS  CONCURRENCY   FILES  DIRECTORIES  OBJECTS         SIZE     MEMORY  EXECUTION TIME  SPEED INCREASE
 WalkDir          N/A      N/A          477859        85515   563374  11650846047  482210408    1.797338872s  1.0x
 Concurrent Read  8        3            477859        85515   563374  11650846047  515456544    849.914464ms  2.1x


 Test Number: 4 | Directory Scan Comparison: /home/matt/projects
 FUNCTION         NUMCPUS  CONCURRENCY   FILES  DIRECTORIES  OBJECTS         SIZE     MEMORY  EXECUTION TIME  SPEED INCREASE
 WalkDir          N/A      N/A          477859        85515   563374  11650846047  482200984    1.792817841s  1.0x
 Concurrent Read  8        4            477859        85515   563374  11650846047  504578416    603.238314ms  3.0x


 Test Number: 5 | Directory Scan Comparison: /home/matt/projects
 FUNCTION         NUMCPUS  CONCURRENCY   FILES  DIRECTORIES  OBJECTS         SIZE     MEMORY  EXECUTION TIME  SPEED INCREASE
 WalkDir          N/A      N/A          477859        85515   563374  11650846047  482151680    1.745606908s  1.0x
 Concurrent Read  8        5            477859        85515   563374  11650846047  520826120     511.29657ms  3.4x


 Test Number: 6 | Directory Scan Comparison: /home/matt/projects
 FUNCTION         NUMCPUS  CONCURRENCY   FILES  DIRECTORIES  OBJECTS         SIZE     MEMORY  EXECUTION TIME  SPEED INCREASE
 WalkDir          N/A      N/A          477859        85515   563374  11650846047  482210200    1.723066673s  1.0x
 Concurrent Read  8        6            477859        85515   563374  11650846047  513077144    430.436392ms  4.0x


 Test Number: 7 | Directory Scan Comparison: /home/matt/projects
 FUNCTION         NUMCPUS  CONCURRENCY   FILES  DIRECTORIES  OBJECTS         SIZE     MEMORY  EXECUTION TIME  SPEED INCREASE
 WalkDir          N/A      N/A          477859        85515   563374  11650846047  482178016    1.775414641s  1.0x
 Concurrent Read  8        7            477859        85515   563374  11650846047  520799128    397.080972ms  4.5x


 Test Number: 8 | Directory Scan Comparison: /home/matt/projects
 FUNCTION         NUMCPUS  CONCURRENCY   FILES  DIRECTORIES  OBJECTS         SIZE     MEMORY  EXECUTION TIME  SPEED INCREASE
 WalkDir          N/A      N/A          477859        85515   563374  11650846047  482200984    1.946071665s  1.0x
 Concurrent Read  8        8            477859        85515   563374  11650846047  520798416    371.251936ms  5.2x

 ```