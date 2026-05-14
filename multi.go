package multicsv

import (
	"errors"
	"io"
)

// Reader is an interface that wraps the basic Read method.
type Reader interface {
	Read() (record []string, err error)
}

type eofReader struct{}

func (eofReader) Read() (record []string, err error) {
	return []string{}, io.EOF
}

// MultiReader is the logical concatenation of the provided input readers.
// They're read sequentially. Once all inputs have returned EOF, Read will return EOF.
// If any of the readers return a non-nil, non-EOF error, Read will return that error.
type MultiReader struct {
	readers []Reader
}

// Read reads one record (a slice of fields) from the provided input readers.
// Following code was taken from https://go.dev/src/io/multi.go and adopted to works with csv readers.
func (mr *MultiReader) Read() (record []string, err error) {
	for len(mr.readers) > 0 {
		if len(mr.readers) == 1 {
			if r, ok := mr.readers[0].(*MultiReader); ok {
				mr.readers = r.readers
				continue
			}
		}

		record, err = mr.readers[0].Read()
		if err == io.EOF {
			closeErr := closeReader(mr.readers[0])
			mr.readers[0] = eofReader{} // permit earlier GC
			mr.readers = mr.readers[1:]

			if closeErr != nil {
				return nil, closeErr
			}
		}

		if len(record) > 0 || err != io.EOF {
			if err == io.EOF && len(mr.readers) > 0 {
				err = nil
			}

			return
		}
	}

	return []string{}, io.EOF
}

// ReadAll reads all the remaining records from the provided input readers.
func (mr *MultiReader) ReadAll() (records [][]string, err error) {
	defer func() {
		if closeErr := mr.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	for {
		record, err := mr.Read()
		if err == io.EOF {
			return records, nil
		}

		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}
}

// NewReader returns a Reader that's the logical concatenation of
// the provided input readers. They're read sequentially. Once all
// inputs have returned EOF, Read will return EOF.  If any of the readers
// return a non-nil, non-EOF error, Read will return that error.
func NewReader(readers ...Reader) *MultiReader {
	r := make([]Reader, len(readers))
	copy(r, readers)

	return &MultiReader{r}
}

// Close closes all remaining readers that implement io.Closer.
func (mr *MultiReader) Close() error {
	var err error

	for _, r := range mr.readers {
		err = errors.Join(err, closeReader(r))
	}

	mr.readers = nil

	return err
}

func closeReader(r Reader) error {
	c, ok := r.(io.Closer)
	if !ok {
		return nil
	}

	return c.Close()
}
