package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var relativePath = ""

func init() {
	filename := os.Args[0] // get command line first parameter
	filedirectory := filepath.Dir(filename)
	relativePath, _ = filepath.Abs(filedirectory)
}

// GetRelativePath implementation
func GetRelativePath() string {
	if relativePath == "" {
		filename := os.Args[0] // get command line first parameter

		filedirectory := filepath.Dir(filename)

		relativePath, _ = filepath.Abs(filedirectory)
		return relativePath
	}
	return relativePath
}

// LoadFile loads a file data in to a byte array
func LoadFile(filepath string) []byte {
	data, err := ioutil.ReadFile(filepath)
	checkErr(err)
	return data
}

// LoadFileRelative loads a file with relative path
func LoadFileRelative(relPath string) []byte {
	filepath := relativePath + "/" + relPath
	return LoadFile(filepath)
}

/* LoadView is a wrapper around LoadFile to load templates
func LoadView(viewName string) []byte {
	filepath := relativePath + "/public/" + viewName + ".html"
	return LoadFile(filepath)
}
*/
