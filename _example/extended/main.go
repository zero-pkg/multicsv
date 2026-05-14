package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/zero-pkg/multicsv"
)

func main() {
	r := multicsv.NewReader(
		customReader("data/count_10.csv"),
		customReader("data/count_100.csv"),
	)
	defer r.Close()

	for {
		line, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}

		fmt.Println(line)
	}
}

func customReader(path string) *multicsv.LazyReader {
	var file *os.File

	return &multicsv.LazyReader{
		Init: func() (*csv.Reader, error) {
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}

			// customize csv.Reader
			r := csv.NewReader(f)
			r.LazyQuotes = true
			file = f

			return r, nil
		},
		CloseFunc: func() error {
			if file == nil {
				return nil
			}

			return file.Close()
		},
	}
}
