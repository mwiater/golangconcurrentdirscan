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
go run . -path="/home/matt/projects"
```

Windows:
```sh
go run . -path="C:\Users\Matt\projects"
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
 Directory Scan Comparison: /home/matt/projects
 FUNCTION         CONCURRENCY  FILES  DIRECTORIES  OBJECTS       SIZE   MEMORY  EXECUTION TIME  SPEED INCREASE 
 WalkDir          N/A           2891         1600     4491  181469808  4611752      1.2886646s  1.0x
 Concurrent Read  1             2891         1600     4491  181469808  5589024      719.8068ms  1.8x


 Directory Scan Comparison: /home/matt/projects
 FUNCTION         CONCURRENCY  FILES  DIRECTORIES  OBJECTS       SIZE   MEMORY  EXECUTION TIME  SPEED INCREASE 
 WalkDir          N/A           2891         1600     4491  181469808  4608032      1.2606058s  1.0x
 Concurrent Read  2             2891         1600     4491  181469808  5102952      448.7249ms  2.8x


 Directory Scan Comparison: /home/matt/projects
 FUNCTION         CONCURRENCY  FILES  DIRECTORIES  OBJECTS       SIZE   MEMORY  EXECUTION TIME  SPEED INCREASE 
 WalkDir          N/A           2891         1600     4491  181469808  4607904      1.2832982s  1.0x
 Concurrent Read  3             2891         1600     4491  181469808  5044456      355.0204ms  3.6x


 Directory Scan Comparison: /home/matt/projects
 FUNCTION         CONCURRENCY  FILES  DIRECTORIES  OBJECTS       SIZE   MEMORY  EXECUTION TIME  SPEED INCREASE 
 WalkDir          N/A           2891         1600     4491  181469808  4607952      1.2155844s  1.0x
 Concurrent Read  4             2891         1600     4491  181469808  4996160      317.6783ms  3.8x


 Directory Scan Comparison: /home/matt/projects
 FUNCTION         CONCURRENCY  FILES  DIRECTORIES  OBJECTS       SIZE   MEMORY  EXECUTION TIME  SPEED INCREASE 
 WalkDir          N/A           2891         1600     4491  181469808  4607952      1.2307747s  1.0x
 Concurrent Read  5             2891         1600     4491  181469808  4988968      342.1745ms  3.6x


 Directory Scan Comparison: /home/matt/projects
 FUNCTION         CONCURRENCY  FILES  DIRECTORIES  OBJECTS       SIZE   MEMORY  EXECUTION TIME  SPEED INCREASE 
 WalkDir          N/A           2891         1600     4491  181469808  4608016      1.2744229s  1.0x
 Concurrent Read  6             2891         1600     4491  181469808  4988664       323.853ms  3.9x

 Directory Scan Comparison: /home/matt/projects
 FUNCTION         CONCURRENCY  FILES  DIRECTORIES  OBJECTS       SIZE   MEMORY  EXECUTION TIME  SPEED INCREASE 
 WalkDir          N/A           2891         1600     4491  181469808  4607920      1.3035374s  1.0x
 Concurrent Read  7             2891         1600     4491  181469808  4993080       355.588ms  3.7x

 Directory Scan Comparison: /home/matt/projects
 FUNCTION         CONCURRENCY  FILES  DIRECTORIES  OBJECTS       SIZE   MEMORY  EXECUTION TIME  SPEED INCREASE 
 WalkDir          N/A           2891         1600     4491  181469808  4608000      1.2724087s  1.0x
 Concurrent Read  8             2891         1600     4491  181469808  4981288      306.3167ms  4.2x
 ```