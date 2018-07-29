package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	targets := flag.Args()
	if len(targets) == 0 {
		targets = []string{"./"}
	}
	fs := RealFileSystem{}
	fmt.Print(SuperProcess(targets, &fs, &options))
}

func isDirectory(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	return file.IsDir()
}

// A DirectoryListing contains all the required data about files in a directory
type DirectoryListing []FileInfo

// A FileInfo contains all the required data about a single file
type FileInfo struct {
	Name  string
	Size  int64
	IsDir bool
}

type Options struct {
	ShowDetails       bool
	SortBySize        bool
	ShowHumanReadable bool
	ReverseSort       bool
}

func (listing DirectoryListing) render(options *Options, buffer *bytes.Buffer) {
	first := true
	if options.ReverseSort {
		for i := len(listing) - 1; i >= 0; i-- {
			renderItem(options, listing[i], first, buffer)
			first = false
		}
	} else {
		for _, file := range listing {
			renderItem(options, file, first, buffer)
			first = false
		}
	}
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
	if size < 1024 {
		return fmt.Sprint(size)
	}
	symbols := []string{"K", "M", "G", "T"}
	symbol := ""
	numbers := float64(size)
	for i := 0; symbol != "T" && numbers > 1024.0; i++ {
		symbol = symbols[i]
		numbers /= 1024
	}
	return fmt.Sprintf("%.1f", numbers) + symbol

}

func (listing DirectoryListing) Process(options *Options, buffer *bytes.Buffer) {
	listing.sort(options)
	listing.render(options, buffer)
}

func (listing DirectoryListing) sort(options *Options) {
	if options.SortBySize {
		sort.SliceStable(listing, func(i, j int) bool { return listing[i].Size > listing[j].Size })
	} else {
		sort.SliceStable(listing, func(i, j int) bool { return strings.ToLower(listing[i].Name) < strings.ToLower(listing[j].Name) })
	}
}

type FileSystem interface {
	ReadDir(string) DirectoryListing
	Stat(string) FileInfo
}

type RealFileSystem struct{}

func (fs *RealFileSystem) Stat(file string) FileInfo {
	info, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}
	return FileInfo{info.Name(), info.Size(), info.IsDir()}

}

func (fs *RealFileSystem) ReadDir(dir string) DirectoryListing {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	result := DirectoryListing{}
	for _, item := range files {
		result = append(result, FileInfo{item.Name(), item.Size(), item.IsDir()})
	}
	return result
}

func SuperProcess(inputs []string, fs FileSystem, options *Options) string {
	buffer := bytes.Buffer{}
	// Print inputs that are files
	for _, input := range inputs {
		var file FileInfo
		file = fs.Stat(input)
		if !file.IsDir {
			DirectoryListing{file}.Process(options, &buffer)
		}
	}

	// Print inputs that are directories
	otherDirs := []string{}
	for _, input := range inputs {
		var file FileInfo
		file = fs.Stat(input)
		if file.IsDir {
			otherDirs = append(otherDirs, input)
		}
	}
	for _, dir := range otherDirs {
		if len(otherDirs) > 1 {
			buffer.WriteString("\n\n" + dir + ":\n\n")
		}
		listing := fs.ReadDir(dir)
		listing.Process(options, &buffer)

	}
	return buffer.String()
}
