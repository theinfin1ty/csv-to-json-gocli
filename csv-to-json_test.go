package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

type FileDataTest struct {
	name    string
	want    InputFile
	wantErr bool
	osArgs  []string
}

type ValidFileTest struct {
	name     string
	filename string
	want     bool
	wantErr  bool
}

type ProcessCsvFileTest struct {
	name      string
	csvString string
	separator string
}

func Test_getFileData(t *testing.T) {
	tests := []FileDataTest{
		{
			name: "Default parameters",
			want: InputFile{
				filePath:  "test.csv",
				separator: ",",
				pretty:    false,
			},
			wantErr: false,
			osArgs:  []string{"cmd", "test.csv"},
		},
		{
			name:    "No parameters",
			want:    InputFile{},
			wantErr: true,
			osArgs:  []string{"cmd"},
		},
		{
			name: "Semicolon enabled",
			want: InputFile{
				filePath:  "test.csv",
				separator: ";",
				pretty:    false,
			},
			wantErr: false,
			osArgs:  []string{"cmd", "--separator=;", "test.csv"},
		},
		{
			name: "Pretty enabled",
			want: InputFile{
				filePath:  "test.csv",
				separator: ",",
				pretty:    true,
			},
			wantErr: false,
			osArgs:  []string{"cmd", "--pretty", "test.csv"},
		},
		{
			name: "Pretty and semicolon enabled",
			want: InputFile{
				filePath:  "test.csv",
				separator: ";",
				pretty:    true,
			},
			wantErr: false,
			osArgs:  []string{"cmd", "--pretty", "--separator=;", "test.csv"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOsArgs := os.Args

			defer func() {
				os.Args = actualOsArgs
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			os.Args = test.osArgs

			got, err := getFileData()

			if err != nil && !test.wantErr {
				t.Errorf("getFileData() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("getFileData() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_checkIfValidFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test.csv")

	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpFile.Name())

	tests := []ValidFileTest{
		{
			name:     "File does exist",
			filename: tmpFile.Name(),
			want:     true,
			wantErr:  false,
		},
		{
			name:     "File does not exist",
			filename: "nowhere/test.csv",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "File is not csv",
			filename: "test.txt",
			want:     false,
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := checkIfValidFile(test.filename)

			if err != nil && !test.wantErr {
				t.Errorf("checkIfValidFile() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if got != test.want {
				t.Errorf("checkIfValidFile() = %v, want %v", got, test.want)
				return
			}
		})
	}
}

func Test_processCsvFile(t *testing.T) {
	wantMapSlice := []map[string]string{
		{"COL1": "1", "COL2": "2", "COL3": "3"},
		{"COL1": "4", "COL2": "5", "COL3": "6"},
	}

	tests := []ProcessCsvFileTest{
		{"Comma separator", "COL1,COL2,COL3\n1,2,3\n4,5,6\n", ","},
		{"Semicolon separator", "COL1;COL2;COL3\n1;2;3\n4;5;6\n", ";"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "test.csv")
			check(err)
			defer os.Remove(tmpfile.Name())

			_, err = tmpfile.WriteString(test.csvString)
			tmpfile.Sync()

			testFileData := InputFile{
				filePath:  tmpfile.Name(),
				pretty:    false,
				separator: test.separator,
			}

			writerChannel := make(chan map[string]string)

			go processCsvFile(testFileData, writerChannel)

			for _, wantMap := range wantMapSlice {
				record := <-writerChannel
				if !reflect.DeepEqual(record, wantMap) {
					t.Errorf("processCsvFile() = %v, want %v", record, wantMap)
					return
				}
			}
		})
	}
}
