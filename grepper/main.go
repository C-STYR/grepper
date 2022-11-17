package main

import (
	"sync"
	
	"grepper/buildTree"
	"grepper/tasklist"
	"grepper/search"

	"github.com/alexflint/go-arg"
)

var args struct {
	SearchTerm string `arg:"positional,required"` // required
	SearchDir  string `arg:"positional"`          //directory to search in, not required, will default to current dir
}

func main() {

	// command line validation tool
	arg.MustParse(&args)

	// waitgroup for search goroutines
	var searchWg sync.WaitGroup

	// tasklist with cap 100
	tl := tasklist.CreateTLChannel(100)

	// searches are returned here
	results := make(chan search.Result, 100)

	// define number of concurrent searchroutines
	searchRoutines := 10

	searchWg.Add(1)

	// NEED METHOD OF STOPPING SEARCH WHEN END OF TASKLIST REACHED



	buildTree.GatherFilenames(".", &tl)
}
