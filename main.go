package main

import (
	"flag"
	"fmt"
	"ls2/internal"
)

func main() {
	options := internal.Options{}
	flag.BoolVar(&options.SortBySize, "S", false, "sort by size")
	flag.BoolVar(&options.ShowDetails, "l", false, "show details")
	flag.BoolVar(&options.ShowHumanReadable, "h", false, "show human readable size")
	flag.BoolVar(&options.ReverseSort, "r", false, "reverse sort order")
	flag.Parse()
	targets := flag.Args()
	if len(targets) == 0 {
		targets = []string{"./"}
	}
	fs := internal.RealFileSystem{}
	fmt.Print(internal.ProcessInputs(targets, &fs, &options))
	fmt.Print("\n")
}
