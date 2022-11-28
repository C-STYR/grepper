package buildTree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	. "grepper/tasklist"
)

// should build and recursively traverse nested paths of the dir/file supplied

func GatherFilenames(path string, tl *Tasklist, wg *sync.WaitGroup) {

	wg.Add(1)
	defer wg.Done()

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("ReadDir error:", err)
		return
	}

	// files is of type []DirEntry and must be parsed for readability
	for _, file := range files {

		// check if file is a dir or file
		if file.IsDir() {

			// exclude the .git directory
			if strings.HasPrefix(file.Name(), ".git") {
				continue

			// else recursively search directory
			} else {
				innerDir := filepath.Join(path, file.Name())
				GatherFilenames(innerDir, tl, wg)
			}
		} else {
			// add the filepath to the tasklist
			newPath := CreateTask(filepath.Join(path, file.Name()))
			tl.Enqueue(newPath)
		}
	}
}
