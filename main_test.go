package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"
)

// Helper function to create temporary test files
func createTestFile(t *testing.T, filename string, lines []string) {
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", filename, err)
	}
	defer file.Close()

	for _, line := range lines {
		file.WriteString(line + "\n")
	}
}

// Helper function to clean up test files
func cleanupTestFiles(filenames ...string) {
	for _, filename := range filenames {
		os.Remove(filename)
	}
}

func TestReadLines(t *testing.T) {
	testFile := "test_read.txt"
	testLines := []string{"line1", "line2", "line3", ""}
	createTestFile(t, testFile, testLines)
	defer cleanupTestFiles(testFile)

	lines, err := readLines(testFile)
	if err != nil {
		t.Fatalf("readLines failed: %v", err)
	}

	// Should skip empty lines
	expected := []string{"line1", "line2", "line3"}
	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("Expected %v, got %v", expected, lines)
	}
}

func TestReadLinesNonExistentFile(t *testing.T) {
	_, err := readLines("nonexistent.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestFindCommonLinesBasic(t *testing.T) {
	files := []FileData{
		{
			Path: "file1.txt",
			Lines: map[string]bool{
				"apple":  true,
				"banana": true,
				"cherry": true,
			},
		},
		{
			Path: "file2.txt",
			Lines: map[string]bool{
				"banana": true,
				"cherry": true,
				"date":   true,
			},
		},
	}

	common := findCommonLines(files)
	expected := []string{"banana", "cherry"}
	sort.Strings(expected)

	if !reflect.DeepEqual(common, expected) {
		t.Errorf("Expected %v, got %v", expected, common)
	}
}

func TestFindCommonLinesThreeFiles(t *testing.T) {
	files := []FileData{
		{
			Path: "file1.txt",
			Lines: map[string]bool{
				"apple":  true,
				"banana": true,
				"cherry": true,
			},
		},
		{
			Path: "file2.txt",
			Lines: map[string]bool{
				"banana": true,
				"cherry": true,
				"date":   true,
			},
		},
		{
			Path: "file3.txt",
			Lines: map[string]bool{
				"banana":     true,
				"elderberry": true,
			},
		},
	}

	common := findCommonLines(files)
	expected := []string{"banana"}

	if !reflect.DeepEqual(common, expected) {
		t.Errorf("Expected %v, got %v", expected, common)
	}
}

func TestFindCommonLinesNoCommon(t *testing.T) {
	files := []FileData{
		{
			Path: "file1.txt",
			Lines: map[string]bool{
				"apple": true,
			},
		},
		{
			Path: "file2.txt",
			Lines: map[string]bool{
				"banana": true,
			},
		},
	}

	common := findCommonLines(files)
	if len(common) != 0 {
		t.Errorf("Expected empty slice, got %v", common)
	}
}

func TestFindCommonLinesEmptyFiles(t *testing.T) {
	files := []FileData{}
	common := findCommonLines(files)
	if common != nil {
		t.Errorf("Expected nil for empty files, got %v", common)
	}
}

func TestFindUniqueLines(t *testing.T) {
	files := []FileData{
		{
			Path: "file1.txt",
			Lines: map[string]bool{
				"apple":  true,
				"banana": true,
				"cherry": true,
			},
		},
		{
			Path: "file2.txt",
			Lines: map[string]bool{
				"banana":     true,
				"date":       true,
				"elderberry": true,
			},
		},
	}

	unique := findUniqueLines(files)

	expectedFile1 := []string{"apple", "cherry"}
	expectedFile2 := []string{"date", "elderberry"}
	sort.Strings(expectedFile1)
	sort.Strings(expectedFile2)

	if !reflect.DeepEqual(unique[0], expectedFile1) {
		t.Errorf("File1 unique: expected %v, got %v", expectedFile1, unique[0])
	}
	if !reflect.DeepEqual(unique[1], expectedFile2) {
		t.Errorf("File2 unique: expected %v, got %v", expectedFile2, unique[1])
	}
}

func TestFindUniqueLinesTotallyDifferent(t *testing.T) {
	files := []FileData{
		{
			Path: "file1.txt",
			Lines: map[string]bool{
				"apple":  true,
				"banana": true,
			},
		},
		{
			Path: "file2.txt",
			Lines: map[string]bool{
				"cherry": true,
				"date":   true,
			},
		},
	}

	unique := findUniqueLines(files)

	// All lines should be unique since no overlap
	if len(unique[0]) != 2 || len(unique[1]) != 2 {
		t.Errorf("Expected 2 unique lines per file, got %d and %d", len(unique[0]), len(unique[1]))
	}
}

func TestFindUniqueLinesTotallyIdentical(t *testing.T) {
	files := []FileData{
		{
			Path: "file1.txt",
			Lines: map[string]bool{
				"apple":  true,
				"banana": true,
			},
		},
		{
			Path: "file2.txt",
			Lines: map[string]bool{
				"apple":  true,
				"banana": true,
			},
		},
	}

	unique := findUniqueLines(files)

	// No lines should be unique since files are identical
	if len(unique[0]) != 0 || len(unique[1]) != 0 {
		t.Errorf("Expected 0 unique lines per file, got %d and %d", len(unique[0]), len(unique[1]))
	}
}

func TestFindPartiallySharedLines(t *testing.T) {
	files := []FileData{
		{
			Path: "file1.txt",
			Lines: map[string]bool{
				"apple":  true,
				"banana": true,
				"cherry": true,
			},
		},
		{
			Path: "file2.txt",
			Lines: map[string]bool{
				"banana": true,
				"cherry": true,
				"date":   true,
			},
		},
		{
			Path: "file3.txt",
			Lines: map[string]bool{
				"cherry":     true,
				"date":       true,
				"elderberry": true,
			},
		},
	}

	partiallyShared := findPartiallySharedLines(files)

	// "banana" appears in files 0,1 (not in file 2)
	// "date" appears in files 1,2 (not in file 0)
	// "cherry" appears in all files, so it shouldn't be in partially shared

	if len(partiallyShared) != 2 {
		t.Errorf("Expected 2 partially shared lines, got %d", len(partiallyShared))
	}

	if indices, exists := partiallyShared["banana"]; !exists {
		t.Error("Expected 'banana' to be partially shared")
	} else if !reflect.DeepEqual(indices, []int{0, 1}) {
		t.Errorf("Expected 'banana' in files [0,1], got %v", indices)
	}

	if indices, exists := partiallyShared["date"]; !exists {
		t.Error("Expected 'date' to be partially shared")
	} else if !reflect.DeepEqual(indices, []int{1, 2}) {
		t.Errorf("Expected 'date' in files [1,2], got %v", indices)
	}
}

func TestIntegrationTwoFiles(t *testing.T) {
	file1 := "test1.txt"
	file2 := "test2.txt"

	file1Lines := []string{"apple", "banana", "cherry", "date"}
	file2Lines := []string{"banana", "cherry", "elderberry", "fig"}

	createTestFile(t, file1, file1Lines)
	createTestFile(t, file2, file2Lines)
	defer cleanupTestFiles(file1, file2)

	// Test reading files
	lines1, err1 := readLines(file1)
	lines2, err2 := readLines(file2)

	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to read test files: %v, %v", err1, err2)
	}

	// Convert to FileData
	files := []FileData{
		{
			Path:  file1,
			Lines: make(map[string]bool),
		},
		{
			Path:  file2,
			Lines: make(map[string]bool),
		},
	}

	for _, line := range lines1 {
		files[0].Lines[line] = true
	}
	for _, line := range lines2 {
		files[1].Lines[line] = true
	}

	// Test common lines
	common := findCommonLines(files)
	expectedCommon := []string{"banana", "cherry"}
	if !reflect.DeepEqual(common, expectedCommon) {
		t.Errorf("Common lines: expected %v, got %v", expectedCommon, common)
	}

	// Test unique lines
	unique := findUniqueLines(files)
	expectedUnique1 := []string{"apple", "date"}
	expectedUnique2 := []string{"elderberry", "fig"}

	if !reflect.DeepEqual(unique[0], expectedUnique1) {
		t.Errorf("File1 unique: expected %v, got %v", expectedUnique1, unique[0])
	}
	if !reflect.DeepEqual(unique[1], expectedUnique2) {
		t.Errorf("File2 unique: expected %v, got %v", expectedUnique2, unique[1])
	}
}

func TestIntegrationThreeFiles(t *testing.T) {
	file1 := "test1.txt"
	file2 := "test2.txt"
	file3 := "test3.txt"

	file1Lines := []string{"apple", "banana", "cherry"}
	file2Lines := []string{"banana", "cherry", "date"}
	file3Lines := []string{"cherry", "date", "elderberry"}

	createTestFile(t, file1, file1Lines)
	createTestFile(t, file2, file2Lines)
	createTestFile(t, file3, file3Lines)
	defer cleanupTestFiles(file1, file2, file3)

	// Create FileData structures
	files := []FileData{
		{Path: file1, Lines: make(map[string]bool)},
		{Path: file2, Lines: make(map[string]bool)},
		{Path: file3, Lines: make(map[string]bool)},
	}

	allLines := [][]string{file1Lines, file2Lines, file3Lines}
	for i, lines := range allLines {
		for _, line := range lines {
			files[i].Lines[line] = true
		}
	}

	// Test common lines (should only be "cherry")
	common := findCommonLines(files)
	expectedCommon := []string{"cherry"}
	if !reflect.DeepEqual(common, expectedCommon) {
		t.Errorf("Common lines: expected %v, got %v", expectedCommon, common)
	}

	// Test unique lines
	unique := findUniqueLines(files)
	expectedUnique1 := []string{"apple"}
	expectedUnique3 := []string{"elderberry"}

	if !reflect.DeepEqual(unique[0], expectedUnique1) {
		t.Errorf("File1 unique: expected %v, got %v", expectedUnique1, unique[0])
	}
	if unique[1] != nil {
		t.Errorf("File2 unique: expected nil, got %v", unique[1])
	}
	if !reflect.DeepEqual(unique[2], expectedUnique3) {
		t.Errorf("File3 unique: expected %v, got %v", expectedUnique3, unique[2])
	}

	// Test partially shared lines
	partiallyShared := findPartiallySharedLines(files)

	// "banana" should be in files 0,1
	// "date" should be in files 1,2
	if len(partiallyShared) != 2 {
		t.Errorf("Expected 2 partially shared lines, got %d", len(partiallyShared))
	}
}

func TestEdgeCaseEmptyFiles(t *testing.T) {
	file1 := "empty1.txt"
	file2 := "empty2.txt"

	createTestFile(t, file1, []string{})
	createTestFile(t, file2, []string{})
	defer cleanupTestFiles(file1, file2)

	files := []FileData{
		{Path: file1, Lines: make(map[string]bool)},
		{Path: file2, Lines: make(map[string]bool)},
	}

	common := findCommonLines(files)
	if len(common) != 0 {
		t.Errorf("Expected no common lines for empty files, got %v", common)
	}

	unique := findUniqueLines(files)
	if len(unique[0]) != 0 || len(unique[1]) != 0 {
		t.Errorf("Expected no unique lines for empty files, got %v", unique)
	}
}

func TestEdgeCaseWhitespaceAndEmptyLines(t *testing.T) {
	file1 := "whitespace1.txt"
	file2 := "whitespace2.txt"

	// Include empty lines and whitespace
	file1Lines := []string{"apple", "", "  banana  ", "cherry", ""}
	file2Lines := []string{"  banana  ", "", "date", "cherry"}

	createTestFile(t, file1, file1Lines)
	createTestFile(t, file2, file2Lines)
	defer cleanupTestFiles(file1, file2)

	lines1, _ := readLines(file1)
	lines2, _ := readLines(file2)

	// readLines should skip empty lines but preserve whitespace
	expectedLines1 := []string{"apple", "  banana  ", "cherry"}
	expectedLines2 := []string{"  banana  ", "date", "cherry"}

	if !reflect.DeepEqual(lines1, expectedLines1) {
		t.Errorf("File1 lines: expected %v, got %v", expectedLines1, lines1)
	}
	if !reflect.DeepEqual(lines2, expectedLines2) {
		t.Errorf("File2 lines: expected %v, got %v", expectedLines2, lines2)
	}
}

// Benchmark tests for performance
func BenchmarkFindCommonLinesSmall(b *testing.B) {
	files := []FileData{
		{Lines: map[string]bool{"a": true, "b": true, "c": true}},
		{Lines: map[string]bool{"b": true, "c": true, "d": true}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findCommonLines(files)
	}
}

func BenchmarkFindCommonLinesLarge(b *testing.B) {
	// Create larger test data
	file1Lines := make(map[string]bool)
	file2Lines := make(map[string]bool)

	for i := 0; i < 1000; i++ {
		file1Lines[fmt.Sprintf("line%d", i)] = true
		if i%2 == 0 {
			file2Lines[fmt.Sprintf("line%d", i)] = true
		}
	}
	for i := 1000; i < 1500; i++ {
		file2Lines[fmt.Sprintf("line%d", i)] = true
	}

	files := []FileData{
		{Lines: file1Lines},
		{Lines: file2Lines},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findCommonLines(files)
	}
}
