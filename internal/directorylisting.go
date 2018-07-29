package internal

import (
	"bytes"
	"sort"
	"strings"
)

// A DirectoryListing contains all the required data about files in a directory
type DirectoryListing []FileInfo

// Process and render a DirectoryListing
func (listing DirectoryListing) ProcessAndRender(options *Options, buffer *bytes.Buffer) {
	listing.sort(options)
	listing.render(options, buffer)
}

// Render all items in a DirectoryListing, according to sort order
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

// Sort a directory listing according to name or size
func (listing DirectoryListing) sort(options *Options) {
	if options.SortBySize {
		sort.SliceStable(listing, func(i, j int) bool { return listing[i].Size > listing[j].Size })
	} else {
		sort.SliceStable(listing, func(i, j int) bool { return strings.ToLower(listing[i].Name) < strings.ToLower(listing[j].Name) })
	}
}
