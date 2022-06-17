package multicsv

import (
	"encoding/csv"
	"os"
	"sync"
)

type LazyReader struct {
	ReaderFunc ReaderFunc
	once       sync.Once
	reader     *csv.Reader
}

type ReaderFunc func() (*csv.Reader, error)

func (r *LazyReader) Read() (record []string, err error) {
	r.once.Do(func() {
		r.reader, err = r.ReaderFunc()
	})

	if err != nil {
		return
	}

	return r.reader.Read()
}

func LazyFileReader(filepath string, skipHeader bool) Reader {
	return &LazyReader{
		ReaderFunc: func() (*csv.Reader, error) {
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
