package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

func main() {
	options := Options{}
	flag.BoolVar(&options.SortBySize, "S", false, "sort by size")
	flag.BoolVar(&options.ShowDetails, "l", false, "show details")
	flag.BoolVar(&options.ShowHumanReadable, "h", false, "show human readable size")
	flag.BoolVar(&options.ReverseSort, "r", false, "reverse sort order")
	flag.Parse()
	dirname := "./"
	if rest := flag.Args(); len(rest) != 0 {
		dirname = rest[0]
	}
	fmt.Printf(GetDirectoryListing(dirname).Process(&options))
}

// A DirectoryListing contains all the required data about files in a directory
type DirectoryListing []FileInfo

// A FileInfo contains all the required data about a single file
type FileInfo struct {
	Name string
	Size int64
}

type Options struct {
	ShowDetails       bool
	SortBySize        bool
	ShowHumanReadable bool
	ReverseSort       bool
}

func (listing DirectoryListing) render(options *Options) string {
	buffer := bytes.Buffer{}
	first := true
	if options.ReverseSort {
		for i := len(listing) - 1; i >= 0; i-- {
			renderItem(options, listing[i], first, &buffer)
			first = false
		}
	} else {
		for _, file := range listing {
			renderItem(options, file, first, &buffer)
			first = false
		}
	}
	return buffer.String()
}

func renderItem(options *Options, file FileInfo, first bool, buffer *bytes.Buffer) {
	if options.ShowDetails {
		renderWithDetails(file, first, options, buffer)
	} else {
		renderSimple(file, first, buffer)
	}
}

func renderSimple(file FileInfo, first bool, buffer *bytes.Buffer) {
	if !first {
		buffer.WriteString("    ")
	}
	buffer.WriteString(file.Name)
}

func renderWithDetails(file FileInfo, first bool, options *Options, buffer *bytes.Buffer) {
	if !first {
		buffer.WriteString("\n")
	}
	buffer.WriteString(file.Name + " " + renderSize(file.Size, options))
}

func renderSize(size int64, options *Options) string {
	if options.ShowHumanReadable {
		return formatHumanReadable(size)
	}
	return fmt.Sprint(size)
}

func formatHumanReadable(size int64) string {
	if size < 1000 {
		return fmt.Sprint(size)
	}
	symbols := []string{"K", "M", "G", "T"}
	symbol := ""
	numbers := float64(size)
	for i := 0; symbol != "T" && numbers > 1000.0; i++ {
		symbol = symbols[i]
		numbers /= 1000
	}
	return fmt.Sprintf("%.1f", numbers) + symbol

}

func (listing DirectoryListing) Process(options *Options) string {
	listing.sort(options)
	return listing.render(options)
}

func (listing DirectoryListing) sort(options *Options) {
	if options.SortBySize {
		sort.SliceStable(listing, func(i, j int) bool { return listing[i].Size > listing[j].Size })
	} else {
		sort.SliceStable(listing, func(i, j int) bool { return strings.ToLower(listing[i].Name) < strings.ToLower(listing[j].Name) })
	}
}

func GetDirectoryListing(dirname string) DirectoryListing {
	listing := DirectoryListing{}
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		listing = append(listing, FileInfo{f.Name(), f.Size()})
	}
	return listing
}
