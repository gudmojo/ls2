package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_directory_with_no_files_should_output_nothing(t *testing.T) {
	listing := DirectoryListing{}
	assert.Equal(t, "", listing.Process(&Options{}))
}

func Test_directory_with_one_file_should_output_its_name(t *testing.T) {
	listing := DirectoryListing{FileInfo{"myFile", 0}}
	assert.Equal(t, "myFile", listing.Process(&Options{}))
}

func Test_directory_with_three_files_should_output_their_names_alphabetically(t *testing.T) {
	listing := DirectoryListing{FileInfo{"myFile3", 0}, FileInfo{"MyFile2", 0}, FileInfo{"myFile1", 0}}
	assert.Equal(t, "myFile1    MyFile2    myFile3", listing.Process(&Options{}))
}

func Test_size_option_should_sort_files_by_size(t *testing.T) {
	options := Options{SortBySize: true}
	listing := DirectoryListing{FileInfo{"a", 3}, FileInfo{"b", 1}, FileInfo{"c", 2}}
	assert.Equal(t, "a    c    b", listing.Process(&options))
}

func Test_details_option_should_output_size(t *testing.T) {
	listing := DirectoryListing{FileInfo{"a", 3}, FileInfo{"b", 1}, FileInfo{"c", 2}}
	assert.Equal(t, `a 3
b 1
c 2`, listing.Process(&Options{ShowDetails: true}))
}

func Test_human_readable_option_K_M_G_T(t *testing.T) {
	options := Options{ShowDetails: true, ShowHumanReadable: true, SortBySize: true, ReverseSort: true}
	listing := DirectoryListing{
		FileInfo{"file", 66},
		FileInfo{"kfile", 6666},
		FileInfo{"mfile", 6666666},
		FileInfo{"gfile", 6666666666},
		FileInfo{"tfile", 6666666666666},
		FileInfo{"tfile2", 6666666666666666}}
	assert.Equal(t, `file 66
kfile 6.7K
mfile 6.7M
gfile 6.7G
tfile 6.7T
tfile2 6666.7T`, listing.Process(&options))
}
