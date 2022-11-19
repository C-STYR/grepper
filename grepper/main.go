package main

import (
	"fmt"
	"sync"
	"time"

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

	// waitgroup for search goroutines
	var searchWg sync.WaitGroup
	var GFwg sync.WaitGroup

	// tasklist with cap 100
	tl := tasklist.CreateTLChannel(100)

	// searches are returned here
	results := make(chan search.Result, 100)
	quit := make(chan int, 1)

	// define number of concurrent searchroutines
	searchParty := 10

	searchWg.Add(1)

// BuildTreeRoutine:
	go func() {
		fmt.Println("treebuilding goroutine spawned")
		// in a goroutine, gather filenames to be parsed and send down tl channel
		buildTree.GatherFilenames(".", &tl, &GFwg)

		// once recursive process is done...
		GFwg.Wait() // this is blocking
		fmt.Println("treebuilding goroutine complete")

		searchWg.Done()

		// send quit message
		quit <- 1
	}()

	time.Sleep(1 * time.Second)
	for i := 0; i < searchParty; i++ {

		// increment searchers wg for each member of search party
		searchWg.Add(1)

		go func() {
			defer searchWg.Done() // schedule decrementation of waitgroup

		SearchLoop:
			for {
				select {
				// if there are tasks in the tasklist channel...
				case task := <-tl.Tasks:
					fmt.Println("this is a task:", task)

					// parse them
					searchResult := search.SearchByLine(string(task), args.SearchTerm)

					// if there's a string match...
					if searchResult != nil {
						fmt.Println("found a result")
						// loop through and send to results channel
						for _, r := range searchResult {
							fmt.Println("r:", r)
							results <- r
						}
					}
				case <-quit:
					break SearchLoop
				}
			}
		}()
	}

	var displayWg sync.WaitGroup

	displayWg.Add(1)
	go func() {
		for {
			select {

			//print results as they come in
			case r := <-results:
				fmt.Printf("%v[%v]:%v\n", r.Path, r.LineNumber, r.Line)
			
			
			default:
				fmt.Println("hit default case")
				if len(results) == 0 {
					fmt.Println("no results brah")
					displayWg.Done()
					return
				} else {
					fmt.Println("something ain't right")
				}
			}
		}
	}()
	displayWg.Wait() //block until all complete
	// time.Sleep(1 * time.Second)
	GFwg.Wait()
	fmt.Println("waiting complete")

	// currently main is completing before all the goroutines are.
}
