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
	ShowDetails bool
	SortBySize  bool
}

func (listing DirectoryListing) render(options *Options) string {
	buffer := bytes.Buffer{}
	for i, file := range listing {
		if options.ShowDetails {
			renderWithDetails(file, i, &buffer)
		} else {
			renderSimple(file, i, &buffer)
		}
	}
	return buffer.String()
}

func renderSimple(file FileInfo, index int, buffer *bytes.Buffer) {
	if index != 0 {
		buffer.WriteString("    ")
	}
	buffer.WriteString(file.Name)
}

func renderWithDetails(file FileInfo, index int, buffer *bytes.Buffer) {
	if index != 0 {
		buffer.WriteString("\n")
	}
	buffer.WriteString(file.Name + " " + fmt.Sprint(file.Size))
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
