package multicsv

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

type mock struct {
	rows    [][]string
	pointer int
}

// Read имплементирует интерфейс совместимый с csv.Reader
func (t *mock) Read() (record []string, err error) {
	if t.pointer >= len(t.rows) {
		return nil, io.EOF
	}

	row := t.rows[t.pointer]
	t.pointer++

	return row, nil
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
