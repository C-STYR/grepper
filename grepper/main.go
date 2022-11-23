package main

import (
	"fmt"
	"sync"

	"grepper/buildTree"
	"grepper/search"
	"grepper/tasklist"

	// "grepper/search"

	"github.com/alexflint/go-arg"
)

var args struct {
	SearchTerm string `arg:"positional,required"` // required
	SearchDir  string `arg:"positional"`          //directory to search in, not required, will default to current dir
}

func main() {

	/*
		Two main goroutines will be running here:
		1. building the tree (compiling list of filenames and adding to tasklist)
		2. checking files for string matches
	*/

	// command line validation tool
	arg.MustParse(&args)

	// WAITGROUPS
	var gatherFilenamesWG sync.WaitGroup
	var searchLinesWG sync.WaitGroup
	var displayResultsWG sync.WaitGroup
	var masterWG sync.WaitGroup

	gatherFilenamesWG.Add(1)
	searchLinesWG.Add(1)
	displayResultsWG.Add(1)
	masterWG.Add(2) // one for each remaining waitgroup

	// tasklist with cap 100
	tl := tasklist.CreateTLChannel(100)

	// CHANNELS
	results := make(chan search.Result, 100)
	gatherFilenamesComplete := make(chan int, 1)
	searchLinesComplete := make(chan int, 1)
	displayResultsComplete := make(chan int, 1)

	// master waitgroup goroutine
	go func() {
		for {
			select {
			case <-searchLinesComplete:
				masterWG.Done()
			case <-displayResultsComplete:
				masterWG.Done()
			}
		}
	}()

	// define number of concurrent searchroutines
	searchParty := 10

	// Compile filenames from the file tree
	go func() {
		fmt.Println("treebuilding goroutine spawned")

		// in a goroutine, gather filenames to be parsed and send down tl channel
		buildTree.GatherFilenames(".", &tl, &gatherFilenamesWG)

		// once recursive process is done...
		gatherFilenamesWG.Wait() // this is blocking
		fmt.Println("treebuilding goroutine complete")

		// send quit message
		gatherFilenamesComplete <- 1
	}()

	// MAIN THREAD
	// spawns a search routine for each member of search party
	// each routine continues to take tasks until no tasks are left
	for i := 0; i < searchParty; i++ {

		searchLinesWG.Add(1)

		go func() {
			defer searchLinesWG.Done()

			SearchLoop:
			for {

				// if there are tasks in the tasklist channel...
				task := tl.Dequeue()

				// parse them
				searchResult := search.SearchByLine(string(task), args.SearchTerm)

				// if there's a string match...
				if searchResult != nil {
					// loop through and send to results channel
					for _, r := range searchResult {
						results <- r
					}
				} else {
					fmt.Println("No hits in file", task)
				}

				// this needs to be tested
				result := <-gatherFilenamesComplete 
				
				if result == 1 && len(tl.Tasks) == 0{
					break SearchLoop
				}
			}
		}()
	}

	go func() {
		searchLinesWG.Wait()
		fmt.Println("***** searchLinesWG DONE: Sending Display Quit Message *****")
		searchLinesComplete <- 1
	}()

	displayResultsWG.Add(1)
	go func() {
		for {
			select {

			//print results as they come in
			case r := <-results:
				fmt.Printf("%v[%v]:%v\n", r.Path, r.LineNumber, r.Line)

			case <-displayResultsComplete:
				if len(results) == 0 {
					fmt.Println("*********** received display quit signal *************")
					displayResultsWG.Done()
					return
				} else {
					continue
				}
			}
		}
	}()
	displayResultsWG.Wait() //block until all complete
	fmt.Println("waiting complete")

	// currently main is completing before all the goroutines are.
}
