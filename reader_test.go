package multicsv

// import (
// 	"io"
// )

// type mock struct {
// 	rows    [][]string
// 	pointer int
// }

// // Read имплементирует интерфейс совместимый с csv.Reader
// func (t *mock) Read() (record []string, err error) {
// 	if t.pointer >= len(t.rows) {
// 		return nil, io.EOF
// 	}

// 	row := t.rows[t.pointer]
// 	t.pointer++

// 	return row, nil
// }
