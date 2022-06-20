package multicsv

import (
	"io"
	"testing"
)

func TestMultiReaderRead(t *testing.T) {
	r := NewReader(
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
	)

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

func TestMultiReaderReadAll(t *testing.T) {
	r := NewReader(
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
	)

	records, err := r.ReadAll()

	equals(t, 4, len(records))
	ok(t, err)
}

func TestMultiReaderEmpty(t *testing.T) {
	r := MultiReader{
		readers: []Reader{},
	}

	records, err := r.ReadAll()

	equals(t, 0, len(records))
	ok(t, err)
}
