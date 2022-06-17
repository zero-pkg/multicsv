package multicsv

import (
	"io"
	"testing"
)

func TestMultiReader(t *testing.T) {
	r := multiReader{
		readers: []Reader{
			&mock{
				rows: [][]string{
					{"a", "b", "c"},
					{"a", "b", "c"},
					{"a", "b", "c"},
				},
			},
			&mock{},
			&mock{
				rows: [][]string{
					{"d", "e", "f"},
				},
			},
		},
	}

	var (
		cnt int
		err error
	)

	for {
		_, err = r.Read()
		if err != nil {
			break
		}

		cnt++
	}

	equals(t, 4, cnt)
	equals(t, err, io.EOF)
}

func TestMultiReaderEmpty(t *testing.T) {
	r := multiReader{
		readers: []Reader{},
	}

	var (
		cnt int
		err error
	)

	for {
		_, err = r.Read()
		if err != nil {
			break
		}

		cnt++
	}

	equals(t, 0, cnt)
	equals(t, err, io.EOF)
}
