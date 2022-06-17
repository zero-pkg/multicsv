package multicsv

import (
	"encoding/csv"
	"os"
	"testing"
)

func TestLazyReader(t *testing.T) {
	r := &LazyReader{
		ReaderFunc: func() (*csv.Reader, error) {
			f, err := os.Open("testdata/custom.csv")
			ok(t, err)

			reader := csv.NewReader(f)
			reader.LazyQuotes = true

			_, err = reader.Read() // skip header
			ok(t, err)

			return reader, nil
		},
	}

	var cnt int

	for {
		fields, err := r.Read()
		if err != nil {
			break
		}

		cnt++

		equals(t, 2, len(fields))
	}

	equals(t, 10, cnt)
}

func TestLazyFileReader(t *testing.T) {
	r := LazyFileReader("testdata/basic.csv", true)

	var cnt int

	for {
		fields, err := r.Read()
		if err != nil {
			break
		}

		cnt++

		equals(t, 6, len(fields))
	}

	equals(t, 5, cnt)
}
