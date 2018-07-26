package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	dirname := "./"
	if len(os.Args) >= 2 {
		dirname = os.Args[1]
	}
	fmt.Printf(GetDirectoryListing(dirname).Process())
}

// A DirectoryListing contains all the required data about files in a directory
type DirectoryListing []FileInfo

// A FileInfo contains all the required data about a single file
type FileInfo struct {
	Name string
}

func (listing DirectoryListing) render() string {
	var buffer bytes.Buffer
	for i := 0; i < len(listing); i++ {
		if i != 0 {
			buffer.WriteString("    ")
		}
		buffer.WriteString(listing[i].Name)
	}
	return buffer.String()
}

func (listing DirectoryListing) Process() string {
	sort.Sort(ByName(listing))
	return listing.render()
}

func GetDirectoryListing(dirname string) DirectoryListing {
	listing := DirectoryListing{}
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		listing = append(listing, FileInfo{f.Name()})
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
