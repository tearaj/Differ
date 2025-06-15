package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
)

type FileData struct {
	Path  string
	Lines map[string]bool
}

func main() {
	var showDiff bool
	flag.BoolVar(&showDiff, "diff", false, "Show different lines instead of common lines")
	flag.BoolVar(&showDiff, "d", false, "Show different lines instead of common lines (shorthand)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <file1> <file2> [file3] [file4] ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Finds common or different lines between multiple files\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s file1.txt file2.txt                    # Show lines common to both files\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s file1.txt file2.txt file3.txt          # Show lines common to all 3 files\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -diff file1.txt file2.txt              # Show lines unique to each file\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -d file1.txt file2.txt file3.txt       # Show lines unique to each of 3 files\n", os.Args[0])
	}
	
	flag.Parse()
	
	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Error: At least 2 files are required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	filePaths := flag.Args()
	
	// Read all files
	var files []FileData
	for _, path := range filePaths {
		lines, err := readLines(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", path, err)
			os.Exit(1)
		}
		
		// Convert slice to map for efficient lookup
		lineMap := make(map[string]bool)
		for _, line := range lines {
			lineMap[line] = true
		}
		
		files = append(files, FileData{
			Path:  path,
			Lines: lineMap,
		})
	}

	if showDiff {
		showDifferentLines(files)
	} else {
		showCommonLines(files)
	}
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines if you want, or comment out this condition
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

func showCommonLines(files []FileData) {
	commonLines := findCommonLines(files)
	
	if len(commonLines) == 0 {
		fmt.Printf("No common lines found across all %d files\n", len(files))
		return
	}

	fmt.Printf("Lines common to all %d files:\n", len(files))
	for i, file := range files {
		if i == len(files)-1 {
			fmt.Printf("  %s\n", file.Path)
		} else {
			fmt.Printf("  %s,\n", file.Path)
		}
	}
	fmt.Printf("\nFound %d common lines:\n\n", len(commonLines))
	
	for _, line := range commonLines {
		fmt.Println(line)
	}
}

func showDifferentLines(files []FileData) {
	uniqueLines := findUniqueLines(files)
	
	totalUnique := 0
	for _, lines := range uniqueLines {
		totalUnique += len(lines)
	}
	
	if totalUnique == 0 {
		fmt.Printf("No unique lines found - all files have identical content\n")
		return
	}

	fmt.Printf("Lines unique to each file (total: %d unique lines):\n\n", totalUnique)
	
	for i, file := range files {
		if len(uniqueLines[i]) > 0 {
			fmt.Printf("Lines only in %s (%d lines):\n", file.Path, len(uniqueLines[i]))
			for _, line := range uniqueLines[i] {
				fmt.Printf("  %s\n", line)
			}
			fmt.Println()
		} else {
			fmt.Printf("No unique lines in %s\n\n", file.Path)
		}
	}
	
	// Also show lines that are shared by some but not all files
	if len(files) > 2 {
		partiallyShared := findPartiallySharedLines(files)
		if len(partiallyShared) > 0 {
			fmt.Printf("Lines shared by some files (but not all):\n")
			for line, fileIndices := range partiallyShared {
				fmt.Printf("  \"%s\" appears in:", line)
				for _, idx := range fileIndices {
					fmt.Printf(" %s", files[idx].Path)
				}
				fmt.Println()
			}
		}
	}
}

func findCommonLines(files []FileData) []string {
	if len(files) == 0 {
		return nil
	}
	
	// Start with lines from the first file
	commonMap := make(map[string]bool)
	for line := range files[0].Lines {
		commonMap[line] = true
	}
	
	// Check each subsequent file and keep only lines that exist in all
	for i := 1; i < len(files); i++ {
		newCommon := make(map[string]bool)
		for line := range commonMap {
			if files[i].Lines[line] {
				newCommon[line] = true
			}
		}
		commonMap = newCommon
	}
	
	// Convert to sorted slice
	var common []string
	for line := range commonMap {
		common = append(common, line)
	}
	sort.Strings(common)
	
	return common
}

func findUniqueLines(files []FileData) [][]string {
	result := make([][]string, len(files))
	
	for i, file := range files {
		var unique []string
		
		for line := range file.Lines {
			isUnique := true
			// Check if this line exists in any other file
			for j, otherFile := range files {
				if i != j && otherFile.Lines[line] {
					isUnique = false
					break
				}
			}
			
			if isUnique {
				unique = append(unique, line)
			}
		}
		
		sort.Strings(unique)
		result[i] = unique
	}
	
	return result
}

func findPartiallySharedLines(files []FileData) map[string][]int {
	// Map from line to list of file indices that contain it
	lineToFiles := make(map[string][]int)
	
	for i, file := range files {
		for line := range file.Lines {
			lineToFiles[line] = append(lineToFiles[line], i)
		}
	}
	
	// Keep only lines that appear in more than 1 file but less than all files
	result := make(map[string][]int)
	for line, fileIndices := range lineToFiles {
		if len(fileIndices) > 1 && len(fileIndices) < len(files) {
			result[line] = fileIndices
		}
	}
	
	return result
}