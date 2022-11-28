package search

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// structure for storing result values
type Result struct {
	Path       string
	LineNumber int
	Line       string
}

func CreateResult(path string, lineNum int, line string) Result {
	return Result{path, lineNum, line}
}

type Results []Result

// searches a file line by line, returns a slice of results
func SearchByLine(filepath string, searchTerm string) Results {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	defer file.Close()

	// we will return this object
	results := make([]Result, 0)
	
	//initialize lines at 1
	lineNum := 1
	scanner := bufio.NewScanner(file)

	// search line by line
	for scanner.Scan() {
		line := scanner.Text()

		// if a match is found, create a result and add it to the results slice
		if strings.Contains(line, searchTerm) {
			newResult := CreateResult(filepath, lineNum, line)
			results = append(results, newResult)
		}

		// advance the line number counter at the end of the loop
		lineNum += 1
	}

	// only return results if there are matches
	if len(results) != 0 {
		return results
	}
	return nil
}
