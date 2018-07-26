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
		sort.Sort(sort.Reverse(BySize(listing)))
	} else {
		sort.Sort(ByName(listing))
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

// ByName implements sort.Interface for []FileInfo based on
// the Name field.
type ByName []FileInfo

func (a ByName) Len() int      { return len(a) }
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool {
	return strings.ToLower(a[i].Name) < strings.ToLower(a[j].Name)
}

// BySize implements sort.Interface for []FileInfo based on
// the Size field.
type BySize []FileInfo

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].Size < a[j].Size }
