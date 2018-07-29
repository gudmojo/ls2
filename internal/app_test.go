package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type FakeFile struct {
	stat    FileInfo
	listing DirectoryListing
}

type FakeFileSystem struct {
	cannedResponses map[string]FakeFile
}

func (fs FakeFileSystem) Stat(filename string) FileInfo {
	return fs.cannedResponses[filename].stat
}

func (fs FakeFileSystem) ReadDir(filename string) DirectoryListing {
	return fs.cannedResponses[filename].listing
}

func Test_directory_with_no_files_should_output_nothing(t *testing.T) {
	inputs := []string{"dira"}
	fileSystem := FakeFileSystem{map[string]FakeFile{
		"dira": {FileInfo{"dira", 0, true}, []FileInfo{}}}}
	assert.Equal(t, ``, ProcessInputs(inputs, fileSystem, &Options{}))
}

func Test_directory_with_one_file_should_output_its_name(t *testing.T) {
	inputs := []string{"dira"}
	fileSystem := FakeFileSystem{map[string]FakeFile{
		"dira": {FileInfo{"dira", 0, true}, []FileInfo{{"myFile", 0, false}}}}}
	assert.Equal(t, `myFile`, ProcessInputs(inputs, fileSystem, &Options{}))
}

func Test_directory_with_three_files_should_output_their_names_alphabetically(t *testing.T) {
	inputs := []string{"dira"}
	fileSystem := FakeFileSystem{map[string]FakeFile{
		"dira": {FileInfo{"dira", 0, true}, []FileInfo{
			{"myFile3", 0, false},
			{"MyFile2", 0, false},
			{"myFile1", 0, false},
		}}}}
	assert.Equal(t, `myFile1    MyFile2    myFile3`, ProcessInputs(inputs, fileSystem, &Options{}))
}

func Test_size_option_should_sort_files_by_size(t *testing.T) {
	inputs := []string{"dira"}
	fileSystem := FakeFileSystem{map[string]FakeFile{
		"dira": {FileInfo{"dira", 0, true}, []FileInfo{
			{"a", 3, false},
			{"b", 1, false},
			{"c", 2, false},
		}}}}
	assert.Equal(t, `a    c    b`, ProcessInputs(inputs, fileSystem, &Options{SortBySize: true}))
}

func Test_details_option_should_output_size(t *testing.T) {
	inputs := []string{"dira"}
	fileSystem := FakeFileSystem{map[string]FakeFile{
		"dira": {FileInfo{"dira", 0, true}, []FileInfo{
			{"a", 3, false},
			{"b", 1, false},
			{"c", 2, false},
		}}}}
	assert.Equal(t, `a 3
b 1
c 2`, ProcessInputs(inputs, fileSystem, &Options{ShowDetails: true}))

}

func Test_human_readable_option_K_M_G_T(t *testing.T) {
	options := Options{ShowDetails: true, ShowHumanReadable: true, SortBySize: true, ReverseSort: true}
	inputs := []string{"dira"}
	fileSystem := FakeFileSystem{map[string]FakeFile{
		"dira": {FileInfo{"dira", 0, true}, []FileInfo{
			{"file", 66, false},
			{"kfile", 6666, false},
			{"mfile", 6666666, false},
			{"gfile", 6666666666, false},
			{"tfile", 6666666666666, false},
			{"tfile2", 6666666666666666, false},
		}}}}
	assert.Equal(t, `file 66
kfile 6.5K
mfile 6.4M
gfile 6.2G
tfile 6.1T
tfile2 6063.3T`, ProcessInputs(inputs, fileSystem, &options))

}

// Todo: fix line break counts to make it pretty

func Test_output_when_two_directories_are_input(t *testing.T) {
	inputs := []string{"dira", "dirb"}
	fileSystem := FakeFileSystem{map[string]FakeFile{
		"dira": {FileInfo{"dira", 0, true}, []FileInfo{}},
		"dirb": {FileInfo{"dirb", 0, true}, []FileInfo{
			{"c", 0, false},
			{"d", 0, true}}}}}
	assert.Equal(t, `

dira:



dirb:

c    d`, ProcessInputs(inputs, fileSystem, &Options{}))

}
