# Multi-File Comparison Tool

A simple CLI tool written in Go that finds common or different lines between multiple files. Perfect for comparing datasets, configuration files, lists, or any line-based text data.

## Features

- **Multi-file support**: Compare 2 or more files simultaneously
- **Two comparison modes**: Find common lines (intersection) or different lines (unique to each file)
- **Smart output limiting**: Shows first 20 lines by default to keep output manageable
- **Flexible display options**: Full output, custom limits, or truncated views
- **Partially shared analysis**: Shows lines shared by some files but not all (3+ files)
- **Clean, organized output**: Clear formatting with counts and file attribution

## Installation

1. Make sure you have Go installed on your system
2. Build the executable:
   ```bash
   go build -o differ main.go
   ```

## Usage

### Basic Syntax
```bash
./differ [flags] <file1> <file2> [file3] [file4] ...
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `-diff` | `-d` | Show different lines instead of common lines |
| `-full` | `-f` | Show full output without truncation |
| `-limit N` | `-l N` | Maximum number of lines to show per section (default: 20) |
| `-help` | `-h` | Show help message |

### Examples

#### Find Common Lines (Default Mode)
```bash
# Find lines common to both files (first 20 lines)
./differ file1.txt file2.txt

# Find lines common to all 3 files
./differ file1.txt file2.txt file3.txt

# Show all common lines without truncation
./differ -full file1.txt file2.txt

# Show first 50 common lines
./differ -limit 50 file1.txt file2.txt
```

#### Find Different Lines
```bash
# Find lines unique to each file (first 20 per file)
./differ -diff file1.txt file2.txt

# Show all unique lines
./differ -diff -full file1.txt file2.txt file3.txt

# Show first 10 unique lines per file
./differ -d -l 10 file1.txt file2.txt
```

## Output Examples

### Common Lines Output
```
Lines common to all 2 files:
  products.txt,
  inventory.txt

Found 25 common lines (showing first 20):

apple
banana
cherry
grape
kiwi
...

... and 5 more lines (use -full or -f to see all)
```

### Different Lines Output
```
Lines unique to each file (total: 15 unique lines):

Lines only in file1.txt (8 lines):
  mango
  papaya
  pineapple
  strawberry
  watermelon
  ... and 3 more lines

Lines only in file2.txt (7 lines):
  blueberry
  cranberry
  elderberry
  ... and 4 more lines

Lines shared by some files (but not all):
  "orange" appears in: file1.txt file3.txt
  "peach" appears in: file2.txt file3.txt

Use -full or -f to see all results, or -limit N to show more lines per section
```

## Use Cases

### Data Analysis
- Compare customer lists from different sources
- Find common entries in multiple datasets
- Identify unique records across databases

### Configuration Management
- Compare configuration files across environments
- Find differences between backup versions
- Identify common settings across multiple configs

### List Processing
- Compare word lists, dictionaries, or vocabularies
- Find overlapping entries in multiple catalogs
- Analyze differences in inventory lists

### Development
- Compare dependency lists
- Find common imports across codebases
- Analyze differences in requirements files

## File Format

The tool works with any text file where each line represents a data entry:
- One entry per line
- Empty lines are automatically skipped
- No special formatting required
- Works with any text encoding

## Performance Notes

- Uses hash maps for efficient line lookups
- Memory usage scales with file size and number of unique lines
- Suitable for files with millions of lines
- Output is automatically sorted alphabetically

## Tips

1. **Start with truncated output**: Use default settings first to get an overview
2. **Use `-full` for complete analysis**: When you need to see everything
3. **Combine with shell tools**: Pipe output to `grep`, `sort`, or `wc` for further processing
4. **Custom limits for large datasets**: Use `-limit` to find the right balance for your terminal

## Examples with Real Data

```bash
# Compare email lists
./differ subscribers.txt newsletter.txt

# Find common dependencies
./differ -full requirements-prod.txt requirements-dev.txt

# Analyze log files (first 100 unique entries per file)
./differ -diff -limit 100 app1.log app2.log app3.log

# Find configuration differences
./differ -d config-staging.ini config-production.ini
```

## Contributing

Feel free to submit issues or pull requests to improve the tool. Some potential enhancements:
- Case-insensitive comparison option
- Regular expression filtering
- CSV column comparison
- Output format options (JSON, CSV)
- Ignore whitespace option

## License

This tool is provided as-is for educational and practical use.