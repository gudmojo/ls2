package internal

import (
	"io/ioutil"
	"log"
	"os"
)

// Implements the FileSystem interface
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
