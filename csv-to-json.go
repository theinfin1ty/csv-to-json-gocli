package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type InputFile struct {
	filePath  string
	separator string
	pretty    bool
}

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

func check(e error) {
	if e != nil {
		exitGracefully(e)
	}
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

func checkIfValidFile(filename string) (bool, error) {
	fileExtension := filepath.Ext(filename)

	if fileExtension != ".csv" {
		return false, errors.New("File extension is not valid")
	}

	_, err := os.Stat(filename)

	if err != nil && os.IsNotExist(err) {
		return false, errors.New("File does not exist")
	}

	return true, nil
}

func processLine(headers []string, dataList []string) (map[string]string, error) {
	if len(dataList) != len(headers) {
		return nil, errors.New("Line doesn't match headers format. Skipping")
	}

	recordMap := make(map[string]string)

	for i, name := range headers {
		recordMap[name] = dataList[i]
	}

	return recordMap, nil
}

func processCsvFile(fileData InputFile, writerChannel chan<- map[string]string) {
	file, err := os.Open(fileData.filePath)

	check(err)

	defer file.Close()

	// var headers, line []string

	reader := csv.NewReader(file)

	if fileData.separator == "" {
		fileData.separator = ","
	}

	reader.Comma = rune(fileData.separator[0])

	headers, err := reader.Read()

	check(err)

	for {
		line, err := reader.Read()

		if err == io.EOF {
			close(writerChannel)
			break
		} else if err != nil {
			exitGracefully(err)
		}

		record, err := processLine(headers, line)

		if err != nil {
			fmt.Printf("Line: %sError: %s\n", line, err)
			continue
		}

		writerChannel <- record
	}
}

func main() {
	fileData, err := getFileData()

	fmt.Println(fileData, err)
}
