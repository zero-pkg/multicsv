package multicsv

import (
	"encoding/csv"
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
// Optionally supports the CSV header skip option.
func LazyFileReader(filepath string, skipHeader bool) Reader {
	return &LazyReader{
		Init: func() (*csv.Reader, error) {
			f, err := os.Open(filepath)
			if err != nil {
				return nil, err
			}

			r := csv.NewReader(f)

			if skipHeader {
				if _, err := r.Read(); err != nil {
					return nil, err
				}
			}

			return r, nil
		},
	}
}
