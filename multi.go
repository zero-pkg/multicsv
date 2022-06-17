package multicsv

import (
	"io"
)

type Reader interface {
	Read() (record []string, err error)
}

type eofReader struct{}

func (eofReader) Read() (record []string, err error) {
	return []string{}, io.EOF
}

type multiReader struct {
	readers []Reader
}

// Read reads one record (a slice of fields) from the provided input readers.
// If the record has an unexpected number of fields,
// Read returns the record along with the error ErrFieldCount.
// Except for that case, Read always returns either a non-nil
// record or a non-nil error, but not both.
// If there is no data left to be read, Read returns nil, io.EOF.
// If ReuseRecord is true, the returned slice may be shared
// between multiple calls to Read.
func (mr *multiReader) Read() (record []string, err error) {
	for len(mr.readers) > 0 {
		if len(mr.readers) == 1 {
			// Optimization to flatten nested multiReaders (Issue 13558).
			if r, ok := mr.readers[0].(*multiReader); ok {
				mr.readers = r.readers
				continue
			}
		}

		record, err = mr.readers[0].Read()
		if err == io.EOF {
			// Use eofReader instead of nil to avoid nil panic
			// after performing flatten (Issue 18232).
			mr.readers[0] = eofReader{} // permit earlier GC
			mr.readers = mr.readers[1:]
		}

		if len(record) > 0 || err != io.EOF {
			if err == io.EOF && len(mr.readers) > 0 {
				// Don't return EOF yet. More readers remain.
				err = nil
			}

			return
		}
	}

	return []string{}, io.EOF
}

// MultiReader returns a Reader that's the logical concatenation of
// the provided input readers. They're read sequentially. Once all
// inputs have returned EOF, Read will return EOF.  If any of the readers
// return a non-nil, non-EOF error, Read will return that error.
func MultiReader(readers ...Reader) Reader {
	r := make([]Reader, len(readers))
	copy(r, readers)

	return &multiReader{r}
}
