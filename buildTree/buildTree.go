package buildTree

import (
	"os"
	"fmt"
	"path/filepath"
	"sync"

	. "grepper/tasklist"
)

// should build and recursively traverse nested paths of the dir/file supplied

var GFwg sync.WaitGroup

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
			innerDir := filepath.Join(path, file.Name())
			GatherFilenames(innerDir, tl)
		} else {
			// add the filepath to the tasklist
			tl.Enqueue(CreateTask(filepath.Join(path, file.Name())))
		}
	}
}
