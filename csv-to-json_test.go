package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

type Test struct {
	name    string
	want    InputFile
	wantErr bool
	osArgs  []string
}

func Test_getFileData(t *testing.T) {
	tests := []Test{
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
