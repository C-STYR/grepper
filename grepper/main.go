package main

import (
	. "grepper/tasklist"
	. "grepper/search"
)

func main() {
	tl := CreateTLChannel(100)
	GatherFilenames(".", &tl)
}
