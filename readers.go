package multicsv

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"sync"
)

// LazyReader allows delayed opening of a resource.
// It can be used to delay opening a resource until the resource is actually read.
type LazyReader struct {
	Init   InitFunc
	once   sync.Once
	reader *csv.Reader
}

// InitFunc is called during the first time reading from LazyReader
type InitFunc func() (*csv.Reader, error)

// Read calls Read func from reader that will be returned by InitFunc.
func (r *LazyReader) Read() (record []string, err error) {
	r.once.Do(func() {
		r.reader, err = r.Init()
	})

	if err != nil {
		return
	}

	return r.reader.Read()
}

// LazyFileReader returns a LazyReader with a predefined InitFunc, which can be used in most cases.
// Optionally supports reader options.
func LazyFileReader(filepath string, opts ...ReaderOption) Reader {
	options := newReaderOptions(opts...)

	return &LazyReader{
		Init: func() (*csv.Reader, error) {
			f, err := os.Open(filepath) //nolint:gosec
			if err != nil {
				return nil, err
			}

			return newCSVReader(f, options)
		},
	}
}

// BytesReader returns a LazyReader for the provided CSV data.
// Optionally supports the same options as LazyFileReader.
func BytesReader(data []byte, opts ...ReaderOption) Reader {
	options := newReaderOptions(opts...)

	return &LazyReader{
		Init: func() (*csv.Reader, error) {
			return newCSVReader(bytes.NewReader(data), options)
		},
	}
}

func newCSVReader(source io.ReadSeeker, options readerOptions) (*csv.Reader, error) {
	r := csv.NewReader(source)
	if options.autoDetectDelimiter {
		delimiter, err := detectDelimiter(source)
		if err != nil {
			return nil, err
		}

		r.Comma = delimiter
	}

	if options.skipHeader {
		if _, err := r.Read(); err != nil {
			return nil, err
		}
	}

	return r, nil
}
