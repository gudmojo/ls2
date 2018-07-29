package internal

import (
	"bytes"
	"fmt"
)

// Start up options
type Options struct {
	ShowDetails       bool
	SortBySize        bool
	ShowHumanReadable bool
	ReverseSort       bool
}

// An abstraction of a file system, so that we can plug in a fake for testing
type FileSystem interface {
	ReadDir(string) DirectoryListing
	Stat(string) FileInfo
}

// A FileInfo contains all the required data about a single file. A bit simpler than os.FileInfo
type FileInfo struct {
	Name  string
	Size  int64
	IsDir bool
}

// Iterates over the files and directories that were provided as inputs to the program and renders them
func ProcessInputs(inputs []string, fs FileSystem, options *Options) string {
	buffer := bytes.Buffer{}

	processFileInputs(inputs, fs, options, &buffer)
	processDirectoryInputs(inputs, fs, &buffer, options)

	return buffer.String()
}

// Iterates over those inputs that are files and renders them
func processFileInputs(inputs []string, fs FileSystem, options *Options, buffer *bytes.Buffer) {
	for _, input := range inputs {
		var file FileInfo
		file = fs.Stat(input)
		if !file.IsDir {
			DirectoryListing{file}.ProcessAndRender(options, buffer)
		}
	}
}

// Iterates over those inputs that are directories and renders them
func processDirectoryInputs(inputs []string, fs FileSystem, buffer *bytes.Buffer, options *Options) {
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
		listing.ProcessAndRender(options, buffer)

	}
}

// Render one file, either with details or without
func renderItem(options *Options, file FileInfo, first bool, buffer *bytes.Buffer) {
	if options.ShowDetails {
		renderWithDetails(file, first, options, buffer)
	} else {
		renderSimple(file, first, buffer)
	}
}

// Render one file without details, including white space to separate items
func renderSimple(file FileInfo, first bool, buffer *bytes.Buffer) {
	if !first {
		buffer.WriteString("    ")
	}
	buffer.WriteString(file.Name)
}

// Render one file with details, including white space to separate items
func renderWithDetails(file FileInfo, first bool, options *Options, buffer *bytes.Buffer) {
	if !first {
		buffer.WriteString("\n")
	}
	buffer.WriteString(file.Name + " " + renderSize(file.Size, options))
}

// Render size, either as human readable or in bytes
func renderSize(size int64, options *Options) string {
	if options.ShowHumanReadable {
		return formatHumanReadable(size)
	}
	return fmt.Sprint(size)
}

// Render human readable sizes, automatically scaled in bytes, kilobytes, megabytes, gigabytes or terabytes
func formatHumanReadable(size int64) string {
	const kb = 1024
	if size < kb {
		return fmt.Sprint(size)
	}
	symbols := []string{"K", "M", "G", "T"}
	symbol := ""
	numbers := float64(size)
	for i := 0; symbol != "T" && numbers > kb; i++ {
		symbol = symbols[i]
		numbers /= kb
	}
	return fmt.Sprintf("%.1f", numbers) + symbol
}
