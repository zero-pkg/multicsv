package multicsv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"sync"
)

var (
	errLazyReaderClosed  = errors.New("multicsv: lazy reader is closed")
	errLazyReaderMissing = errors.New("multicsv: lazy reader init returned nil reader")
)

// LazyReader allows delayed opening of a resource.
// It can be used to delay opening a resource until the resource is actually read.
type LazyReader struct {
	Init InitFunc
	// CloseFunc closes resources opened by Init. It is optional and called at most once.
	CloseFunc func() error
	once      sync.Once
	mu        sync.Mutex
	reader    *csv.Reader
	initErr   error
	closed    bool
	eof       bool
}

// InitFunc is called during the first time reading from LazyReader
type InitFunc func() (*csv.Reader, error)

// LazyFileReader returns a LazyReader for the provided CSV file.
// Optionally supports reader options.
func LazyFileReader(filepath string, opts ...ReaderOption) Reader {
	options := newReaderOptions(opts...)

	var file *os.File

	return &LazyReader{
		Init: func() (*csv.Reader, error) {
			f, err := os.Open(filepath) //nolint:gosec
			if err != nil {
				return nil, err
			}

			r, err := newCSVReader(f, options)
			if err != nil {
				return nil, errors.Join(err, f.Close())
			}

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

// Read calls Read func from reader that will be returned by InitFunc.
func (r *LazyReader) Read() (record []string, err error) {
	r.once.Do(func() {
		r.mu.Lock()
		defer r.mu.Unlock()

		if r.closed {
			r.initErr = errLazyReaderClosed
			return
		}

		r.reader, r.initErr = r.Init()
	})

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.initErr != nil {
		return nil, r.initErr
	}

	if r.closed {
		if r.eof {
			return nil, io.EOF
		}

		return nil, errLazyReaderClosed
	}

	if r.reader == nil {
		return nil, errLazyReaderMissing
	}

	record, err = r.reader.Read()
	if err == io.EOF {
		r.eof = true
		if closeErr := r.closeLocked(); closeErr != nil {
			return record, closeErr
		}
	}

	return record, err
}

// Close closes the resource opened by LazyReader.
func (r *LazyReader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.closeLocked()
}

func (r *LazyReader) closeLocked() error {
	if r.closed {
		return nil
	}

	r.closed = true
	if r.CloseFunc == nil {
		return nil
	}

	return r.CloseFunc()
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
