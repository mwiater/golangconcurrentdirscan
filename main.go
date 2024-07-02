package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// ScannedFile holds the path and information about a file or directory.
type ScannedFile struct {
	Path string      // Path is the location of the file or directory.
	Info os.FileInfo // Info is the file information.
}

// ScanResult holds a collection of scanned files and directories, providing thread-safe access.
type ScanResult struct {
	Files []ScannedFile
	sync.Mutex
}

// AddFile safely adds a ScannedFile to the ScanResult.
func (result *ScanResult) AddFile(file ScannedFile) {
	result.Lock()
	defer result.Unlock()
	result.Files = append(result.Files, file)
}

// scanUsingWalkDir scans the directory tree rooted at rootDir using filepath.WalkDir.
// It returns a slice of ScannedFile and any error encountered.
func scanUsingWalkDir(rootDir string) ([]ScannedFile, error) {
	var result ScanResult

	// Walk the directory tree rooted at rootDir.
	err := filepath.WalkDir(rootDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err // If there is an error accessing the path, return the error.
		}
		info, err := entry.Info()
		if err != nil {
			return err // If there is an error retrieving the file info, return the error.
		}
		// Add the file to the scan result.
		result.AddFile(ScannedFile{Path: path, Info: info})
		return nil
	})
	// Return the scanned files and any error encountered during the walk.
	return result.Files, err
}

// scanUsingGoroutines scans the directory tree rooted at rootDir using concurrent goroutines.
// It returns a slice of ScannedFile and any error encountered.
func scanUsingGoroutines(rootDir string, numCPU int, concurrencyMultiplier int) ([]ScannedFile, error) {
	var result ScanResult
	var wg sync.WaitGroup

	// Baseline
	if numCPU == 0 {
		concurrencyMultiplier = 1
		numCPU = 1
	}

	semaphore := make(chan struct{}, numCPU*concurrencyMultiplier) // Buffered channel to limit concurrent goroutines.
	errChan := make(chan error, 1)                                 // Channel to capture errors.
	doneChan := make(chan struct{})                                // Channel to signal completion.

	// scanDirectory is a recursive function to scan directories.
	var scanDirectory func(string)
	scanDirectory = func(dirPath string) {
		defer wg.Done()                // Decrement the WaitGroup counter when the function returns.
		semaphore <- struct{}{}        // Acquire a slot in the semaphore.
		defer func() { <-semaphore }() // Release the slot in the semaphore when done.

		entries, err := os.ReadDir(dirPath)
		if err != nil {
			select {
			case errChan <- err: // Send error to errChan if possible.
			default:
			}
			return
		}

		buffer := make([]ScannedFile, 0, len(entries)) // Buffer to store scanned files.
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				select {
				case errChan <- err: // Send error to errChan if possible.
				default:
				}
				continue
			}
			// Append the scanned file to the buffer.
			buffer = append(buffer, ScannedFile{Path: filepath.Join(dirPath, entry.Name()), Info: info})

			if entry.IsDir() {
				wg.Add(1)                                              // Increment the WaitGroup counter for the new goroutine.
				go scanDirectory(filepath.Join(dirPath, entry.Name())) // Scan the subdirectory.
			}
		}

		// Append the buffer to the result safely.
		result.Lock()
		result.Files = append(result.Files, buffer...)
		result.Unlock()
	}

	// Get information about the root directory.
	rootInfo, err := os.Stat(rootDir)
	if err != nil {
		return nil, err
	}
	result.AddFile(ScannedFile{Path: rootDir, Info: rootInfo}) // Add the root directory to the result.

	wg.Add(1)                 // Increment the WaitGroup counter for the initial scan.
	go scanDirectory(rootDir) // Start scanning from the root directory.

	// Wait for all goroutines to finish and close doneChan.
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	// Wait for completion or an error.
	select {
	case <-doneChan:
		return result.Files, nil // Return the result if done.
	case err := <-errChan:
		return nil, err // Return the error if encountered.
	}
}

// ScanComparison holds the comparison results of two directory scans.
type ScanComparison struct {
	TotalWalkDirFiles     int      // Total number of files found using WalkDir.
	TotalWalkDirDirs      int      // Total number of directories found using WalkDir.
	TotalWalkDirSize      int64    // Total size of all files found using WalkDir.
	TotalWalkDirObjects   int      // Total number of files and directories found using WalkDir.
	TotalGoroutineFiles   int      // Total number of files found using Goroutines.
	TotalGoroutineDirs    int      // Total number of directories found using Goroutines.
	TotalGoroutineSize    int64    // Total size of all files found using Goroutines.
	TotalGoroutineObjects int      // Total number of files and directories found using Goroutines.
	OnlyInWalkDir         []string // Paths only found in WalkDir scan.
	OnlyInGoroutines      []string // Paths only found in Goroutines scan.
}

// compareScanResults compares the results of two directory scans.
// It returns a ScanComparison struct with detailed comparison results.
func compareScanResults(walkDirFiles, goroutineFiles []ScannedFile) ScanComparison {
	walkDirMap := make(map[string]os.FileInfo)
	goroutineMap := make(map[string]os.FileInfo)

	var comparison ScanComparison

	for _, file := range walkDirFiles {
		walkDirMap[file.Path] = file.Info
		if file.Info.IsDir() {
			comparison.TotalWalkDirDirs++
		} else {
			comparison.TotalWalkDirFiles++
			comparison.TotalWalkDirSize += file.Info.Size()
		}
	}
	comparison.TotalWalkDirObjects = comparison.TotalWalkDirFiles + comparison.TotalWalkDirDirs

	for _, file := range goroutineFiles {
		goroutineMap[file.Path] = file.Info
		if file.Info.IsDir() {
			comparison.TotalGoroutineDirs++
		} else {
			comparison.TotalGoroutineFiles++
			comparison.TotalGoroutineSize += file.Info.Size()
		}
	}
	comparison.TotalGoroutineObjects = comparison.TotalGoroutineFiles + comparison.TotalGoroutineDirs

	for path := range walkDirMap {
		if _, found := goroutineMap[path]; !found {
			comparison.OnlyInWalkDir = append(comparison.OnlyInWalkDir, path)
		}
	}
	for path := range goroutineMap {
		if _, found := walkDirMap[path]; !found {
			comparison.OnlyInGoroutines = append(comparison.OnlyInGoroutines, path)
		}
	}

	return comparison
}

// main is the entry point of the program. It scans a directory using both WalkDir and goroutines,
// then compares the results and prints a summary.
func main() {
	var inputPath string
	var concurrencyMultiplier int

	// Define the flags
	flag.StringVar(&inputPath, "path", "", "The path to the projects directory")
	flag.IntVar(&concurrencyMultiplier, "concurrencyMultiplier", 1, "increase default concurrency from runtime.NumCPU() * 2 to runtime.NumCPU() * concurrencyMultiplier")

	// Parse the flags
	flag.Parse()

	// Check if the input path is provided
	if inputPath == "" {
		fmt.Println("Please provide a path using the -path flag.")
		return
	}

	var rootDir string
	if runtime.GOOS == "windows" {
		rootDir = filepath.Clean(inputPath)
	} else {
		rootDir = filepath.Clean(inputPath)
	}

	// Capture memory stats before and after each scan
	var memStatsBefore, memStatsAfter runtime.MemStats

	numCPU := runtime.NumCPU()
	for concurrency := 0; concurrency <= numCPU; concurrency++ {

		// Measure WalkDir performance
		runtime.GC() // Force garbage collection
		runtime.ReadMemStats(&memStatsBefore)
		start := time.Now()
		walkDirFiles, err := scanUsingWalkDir(rootDir)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		durationWalkDir := time.Since(start)
		runtime.ReadMemStats(&memStatsAfter)
		memUsedWalkDir := memStatsAfter.TotalAlloc - memStatsBefore.TotalAlloc

		// Measure Goroutines performance
		runtime.GC() // Force garbage collection
		runtime.ReadMemStats(&memStatsBefore)
		start = time.Now()
		goroutineFiles, err := scanUsingGoroutines(rootDir, concurrency, concurrencyMultiplier)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		durationGoroutines := time.Since(start)
		runtime.ReadMemStats(&memStatsAfter)
		memUsedGoroutines := memStatsAfter.TotalAlloc - memStatsBefore.TotalAlloc

		comparison := compareScanResults(walkDirFiles, goroutineFiles)
		if len(comparison.OnlyInWalkDir) > 0 || len(comparison.OnlyInGoroutines) > 0 {
			fmt.Println("Mismatch in number of files found!")
		}

		var speedIncreaseWalkDirRounded float64
		var speedIncreaseGoroutinesRounded float64

		// Determine the slower and faster execution times
		var speedIncreaseWalkDir, speedIncreaseGoroutines float64
		if durationWalkDir > durationGoroutines {
			speedIncreaseWalkDir = 1.0
			speedIncreaseGoroutines = float64(durationWalkDir) / float64(durationGoroutines)
			speedIncreaseWalkDirRounded = float64(int(speedIncreaseWalkDir*10+0.5)) / 10
			speedIncreaseGoroutinesRounded = float64(int(speedIncreaseGoroutines*10+0.5)) / 10
		} else {
			speedIncreaseWalkDir = float64(durationGoroutines) / float64(durationWalkDir)
			speedIncreaseGoroutines = 1.0
			speedIncreaseWalkDirRounded = float64(int(speedIncreaseWalkDir*10+0.5)) / 10
			speedIncreaseGoroutinesRounded = float64(int(speedIncreaseGoroutines*10+0.5)) / 10
		}

		chartTitle := "Test Number: " + strconv.Itoa(concurrency) + " | Directory Scan Comparison: " + rootDir
		chartConcurrency := strconv.Itoa(concurrency * concurrencyMultiplier)

		if concurrency*concurrencyMultiplier == 0 {
			chartTitle = "Baseline: No Concurrency | Directory Scan Comparison: " + rootDir
			chartConcurrency = "N/A"
		}

		resultsTable := Table("DarkSimple", chartTitle)

		resultsTable.AppendHeader([]interface{}{
			"Function",
			"NumCPUs",
			"Concurrency",
			"Files",
			"Directories",
			"Objects",
			"Size",
			"Memory",
			"Execution Time",
			"Speed Increase",
		})

		resultsTable.AppendRow([]interface{}{
			"WalkDir",
			"N/A",
			"N/A",
			comparison.TotalWalkDirFiles,
			comparison.TotalWalkDirDirs,
			comparison.TotalWalkDirObjects,
			comparison.TotalWalkDirSize,
			memUsedWalkDir,
			durationWalkDir,
			fmt.Sprintf("%.1fx", speedIncreaseWalkDirRounded),
		})

		resultsTable.AppendRow([]interface{}{
			"Concurrent Read",
			numCPU,
			chartConcurrency,
			comparison.TotalGoroutineFiles,
			comparison.TotalGoroutineDirs,
			comparison.TotalGoroutineObjects,
			comparison.TotalGoroutineSize,
			memUsedGoroutines,
			durationGoroutines,
			fmt.Sprintf("%.1fx", speedIncreaseGoroutinesRounded),
		})

		fmt.Println()
		resultsTable.Render()
		fmt.Println()

		if len(comparison.OnlyInWalkDir) > 0 {
			fmt.Println("Files only in WalkDir:")
			for _, path := range comparison.OnlyInWalkDir {
				fmt.Println(path)
			}
		}
		if len(comparison.OnlyInGoroutines) > 0 {
			fmt.Println("Files only in Goroutines:")
			for _, path := range comparison.OnlyInGoroutines {
				fmt.Println(path)
			}
		}
	}
}
