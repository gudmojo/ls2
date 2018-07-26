package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_directory_with_no_files_should_output_nothing(t *testing.T) {
	listing := DirectoryListing{}
	assert.Equal(t, listing.Process(), "")
}

func Test_directory_with_one_file_should_output_its_name(t *testing.T) {
	listing := DirectoryListing{FileInfo{"myFile"}}
	assert.Equal(t, listing.Process(), "myFile")

}

func Test_directory_with_three_files_should_output_their_names_alphabetically(t *testing.T) {
	listing := DirectoryListing{FileInfo{"myFile3"}, FileInfo{"MyFile2"}, FileInfo{"myFile1"}}
	assert.Equal(t, listing.Process(), "myFile1    MyFile2    myFile3")

}
