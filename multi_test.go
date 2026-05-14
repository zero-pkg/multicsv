package multicsv

import (
	"errors"
	"io"
	"testing"
)

type closeMock struct {
	mock
	closeErr   error
	closeCount int
}

func (m *closeMock) Close() error {
	m.closeCount++

	return m.closeErr
}

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

func TestMultiReaderClosesReadersOnEOF(t *testing.T) {
	first := &closeMock{
		mock: mock{
			rows: [][]string{
				{"a", "b", "c"},
			},
		},
	}
	second := &closeMock{
		mock: mock{
			rows: [][]string{
				{"d", "e", "f"},
			},
		},
	}

	r := NewReader(first, second)
	records, err := r.ReadAll()

	ok(t, err)
	equals(t, 2, len(records))
	equals(t, 1, first.closeCount)
	equals(t, 1, second.closeCount)
}

func TestMultiReaderCloseClosesRemainingReaders(t *testing.T) {
	first := &closeMock{
		mock: mock{
			rows: [][]string{
				{"a", "b", "c"},
				{"d", "e", "f"},
			},
		},
	}
	second := &closeMock{
		mock: mock{
			rows: [][]string{
				{"g", "h", "i"},
			},
		},
	}

	r := NewReader(first, second)
	_, err := r.Read()
	ok(t, err)

	err = r.Close()
	ok(t, err)

	equals(t, 1, first.closeCount)
	equals(t, 1, second.closeCount)

	_, err = r.Read()
	equals(t, io.EOF, err)
}

func TestMultiReaderReturnsCloseError(t *testing.T) {
	expectedErr := errors.New("close failed")
	reader := &closeMock{
		mock: mock{
			rows: [][]string{
				{"a", "b", "c"},
			},
		},
		closeErr: expectedErr,
	}

	r := NewReader(reader)
	_, err := r.Read()
	ok(t, err)

	_, err = r.Read()

	assert(t, errors.Is(err, expectedErr), "expected close error, got %v", err)
	equals(t, 1, reader.closeCount)
}
