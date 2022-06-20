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

func customReader(file string) *multicsv.LazyReader {
	return &multicsv.LazyReader{
		Init: func() (*csv.Reader, error) {
			f, err := os.Open(file)
			if err != nil {
				return nil, err
			}

			// customize csv.Reader
			r := csv.NewReader(f)
			r.LazyQuotes = true

			return r, nil
		},
	}
}
