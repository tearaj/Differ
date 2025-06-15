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

type DiffViewerConfig struct {
	ShowDiff  bool
	ShowFull  bool
	MaxLines  int
}
var config = DiffViewerConfig{}
func main() {
	
	flag.BoolVar(&config.ShowDiff, "diff", false, "Show different lines instead of common lines")
	flag.BoolVar(&config.ShowDiff, "d", false, "Show different lines instead of common lines (shorthand)")
	flag.BoolVar(&config.ShowFull, "full", false, "Show full output without truncation")
	flag.BoolVar(&config.ShowFull, "f", false, "Show full output without truncation (shorthand)")
	flag.IntVar(&config.MaxLines, "limit", 20, "Maximum number of lines to show per section (use with -full to override)")
	flag.IntVar(&config.MaxLines, "l", 20, "Maximum number of lines to show per section (shorthand)")
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <file1> <file2> [file3] [file4] ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Finds common or different lines between multiple files\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s file1.txt file2.txt                    # Show first 20 lines common to both files\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -full file1.txt file2.txt              # Show all common lines\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -limit 50 file1.txt file2.txt          # Show first 50 lines\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -diff file1.txt file2.txt              # Show first 20 unique lines per file\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -d -f file1.txt file2.txt              # Show all unique lines\n", os.Args[0])
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

	if config.ShowDiff {
		showDifferentLines(files, config)
	} else {
		showCommonLines(files, config)
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

func showCommonLines(files []FileData, config DiffViewerConfig) {
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
	fmt.Printf("\nFound %d common lines", len(commonLines))
	
	displayLines := commonLines
	if !config.ShowFull && len(commonLines) > config.MaxLines {
		displayLines = commonLines[:config.MaxLines]
		fmt.Printf(" (showing first %d):\n\n", config.MaxLines)
	} else {
		fmt.Printf(":\n\n")
	}
	
	for _, line := range displayLines {
		fmt.Println(line)
	}
	
	if !config.ShowFull && len(commonLines) > config.MaxLines {
		remaining := len(commonLines) - config.MaxLines
		fmt.Printf("\n... and %d more lines (use -full or -f to see all)\n", remaining)
	}
}

func showDifferentLines(files []FileData, config DiffViewerConfig) {
	
	uniqueLines := findUniqueLines(files)
	
	totalUnique := countLines(uniqueLines)
	
	if totalUnique == 0 {
		fmt.Printf("No unique lines found - all files have identical content\n")
		return
	}
	fmt.Printf("Lines unique to each file (total: %d unique lines):\n\n", totalUnique)
	
	for i, file := range files {
		if len(uniqueLines[i]) > 0 {
			newFunction(file, uniqueLines, i, config.ShowFull, config.MaxLines)
		} else {
			fmt.Printf("No unique lines in %s\n\n", file.Path)
		}
	}
	
	// Also show lines that are shared by some but not all files
	if len(files) > 2 {
		partiallyShared := findPartiallySharedLines(files)
		if len(partiallyShared) > 0 {
			fmt.Printf("Lines shared by some files (but not all):\n")
			
			// Convert to slice for easier handling
			var partialLines []string
			for line := range partiallyShared {
				partialLines = append(partialLines, line)
			}
			sort.Strings(partialLines)
			
			displayPartialLines := partialLines
			if !config.ShowFull && len(partialLines) > config.MaxLines {
				displayPartialLines = partialLines[:config.MaxLines]
				fmt.Printf("(showing first %d of %d):\n", config.MaxLines, len(partialLines))
			}
			
			for _, line := range displayPartialLines {
				fileIndices := partiallyShared[line]
				fmt.Printf("  \"%s\" appears in:", line)
				for _, idx := range fileIndices {
					fmt.Printf(" %s", files[idx].Path)
				}
				fmt.Println()
			}
			
			if !config.ShowFull && len(partialLines) > config.MaxLines {
				remaining := len(partialLines) - config.MaxLines
				fmt.Printf("  ... and %d more partially shared lines\n", remaining)
			}
		}
	}
	
	if !config.ShowFull && (totalUnique > config.MaxLines*len(files) || (len(files) > 2 && len(findPartiallySharedLines(files)) > config.MaxLines)) {
		fmt.Printf("\nUse -full or -f to see all results, or -limit N to show more lines per section\n")
	}
}

func newFunction(file FileData, uniqueLines [][]string, i int, showFull bool, maxLines int) {
	fmt.Printf("Lines only in %s (%d lines", file.Path, len(uniqueLines[i]))

	displayLines := uniqueLines[i]
	if !showFull && len(uniqueLines[i]) > maxLines {
		displayLines = uniqueLines[i][:maxLines]
		fmt.Printf(", showing first %d):\n", maxLines)
	} else {
		fmt.Printf("):\n")
	}

	for _, line := range displayLines {
		fmt.Printf("  %s\n", line)
	}

	if !showFull && len(uniqueLines[i]) > maxLines {
		remaining := len(uniqueLines[i]) - maxLines
		fmt.Printf("  ... and %d more lines\n", remaining)
	}
	fmt.Println()
}

func countLines(uniqueLines [][]string) int {
	totalUnique := 0
	for _, lines := range uniqueLines {
		totalUnique += len(lines)
	}
	return totalUnique
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