package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/zero-pkg/multicsv"
)

func main() {
	r, err := NewMultiReader([]string{"data/count_10.csv", "data/count_100.csv", "data/count_1000.csv"}, ',')
	if err != nil {
		panic(err)
	}

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

func NewMultiReader(files []string, comma rune) (multicsv.Reader, error) {
	readers := make([]multicsv.Reader, len(files))
	for i := range files {
		i := i
		readers[i] = &multicsv.LazyReader{
			Init: func() (*csv.Reader, error) {
				f, err := os.Open(files[i])
				if err != nil {
					return nil, err
				}

				r := csv.NewReader(f)
				r.LazyQuotes = true

				return r, nil
			},
		}
	}

	return multicsv.MultiReader(readers...), nil
}
