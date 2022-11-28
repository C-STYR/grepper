package main

import (
	"fmt"
	"sync"

	"grepper/buildTree"
	"grepper/search"
	"grepper/tasklist"

	"github.com/alexflint/go-arg"
)

var args struct {
	SearchTerm string `arg:"positional,required"` // required
	SearchDir  string `arg:"positional"`          // directory to search in, not required, will default to current dir
}

func main() {

	// command line validation tool
	arg.MustParse(&args)

	// WAITGROUPS
	var gatherFilenamesWG sync.WaitGroup
	var searchLinesWG sync.WaitGroup
	var masterWG sync.WaitGroup

	// one for each process: gathering paths, searching files, displaying results
	masterWG.Add(3) 

	// tasklist with cap 100
	tl := tasklist.CreateTLChannel(100)

	// CHANNELS
	results := make(chan search.Result, 100) // to send search results to display process
	gatherFilenamesComplete := make(chan int) // to inform search routine no more incoming paths
	searchLinesComplete := make(chan int) // to inform display routine no more incoming results

	// Compile filenames from the file tree
	go func() {
		fmt.Println("----- Gathering Paths -----")

		// in a goroutine, gather filenames to be parsed and send down tl channel
		buildTree.GatherFilenames(".", &tl, &gatherFilenamesWG)

		gatherFilenamesWG.Wait() // wait for recursive search to complete

		close(gatherFilenamesComplete) // closed channel will signify process complete
		
		masterWG.Done() // first master process complete
	}()

	/* MAIN THREAD
		- spawns a search routine for each member of search party
		- each routine continues to take tasks until no tasks are left
	*/
	
	searchParty := 10 // define number of concurrent searchroutines
	fmt.Println("----- Starting Search -----")

	for i := 0; i < searchParty; i++ {

		searchLinesWG.Add(1)

		go func() {
			defer searchLinesWG.Done()

		SearchLoop:
			for {
				select {

				// if there are tasks in the tasklist channel...
				case task := <-tl.Tasks:

					// parse them
					searchResult := search.SearchByLine(string(task), args.SearchTerm)

						// loop through and send to results channel
						for _, r := range searchResult {
							results <- r
						}

				// if the gFC channel is closed (which signifies that process is complete)...
				case <-gatherFilenamesComplete:

					// check if there are filenames remaining in the tasks channel
					if len(tl.Tasks) == 0 {
						break SearchLoop
					}
				}
			}
		}()
	}

	// this routine waits until all search routines finish, then closes channel
	go func() {
		searchLinesWG.Wait()
		close(searchLinesComplete)
		masterWG.Done() // second master process complete
	}()

	// this routine displays results as they arrive in the results channel
	go func() {
		for {
			select {
			case r := <-results:
				fmt.Printf("%v[%v]:%v\n", r.Path, r.LineNumber, r.Line)

			// if the search channel is closed, check for results waiting to be displayed
			case <-searchLinesComplete:
				if len(results) == 0 {
					masterWG.Done() // third master process complete
					return
				} else {
					continue
				}
			}
		}
	}()
	masterWG.Wait() //blocks main thread until all processes are complete
	fmt.Println("----- Search Complete -----")
}
