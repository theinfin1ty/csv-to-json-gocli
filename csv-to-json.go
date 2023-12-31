package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type InputFile struct {
	filePath  string
	separator string
	pretty    bool
}

func getFileData() (InputFile, error) {
	if len(os.Args) < 2 {
		return InputFile{}, errors.New("Filepath argument is missing")
	}

	separator := flag.String("separator", ",", "Column separator")
	pretty := flag.Bool("pretty", false, "Generate pretty JSON")

	flag.Parse()

	fileLocation := flag.Arg(0)

	return InputFile{
		filePath:  fileLocation,
		separator: *separator,
		pretty:    *pretty,
	}, nil
}

func main() {
	fileData, err := getFileData()

	fmt.Println(fileData, err)
}
