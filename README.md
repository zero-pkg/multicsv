# multicsv

[![build](https://github.com/zero-pkg/multicsv/actions/workflows/ci.yml/badge.svg)](https://github.com/zero-pkg/multicsv/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/zero-pkg/multicsv)](https://goreportcard.com/report/github.com/zero-pkg/multicsv)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/zero-pkg/multicsv/blob/master/LICENSE)

MultiCSV is a multi reader that's the logical concatenation of the provided input csv readers.
They're read sequentially. Once all inputs have returned EOF, Read will return EOF.
If any of the readers return a non-nil, non-EOF error, Read will return that error.

## Install and update

`go get -u github.com/zero-pkg/multicsv`

## How to use

```go
r := multicsv.NewReader(
    multicsv.LazyFileReader("data/users.csv", true),
    multicsv.LazyFileReader("data/users2.csv", true),
)

records, err := r.ReadAll()
if err != nil {
    panic(err)
}

fmt.Println(records)
```

## Extending LazyReader

```go
func main() {
	r := multicsv.NewReader(
		customReader("data/count_10.csv"),
		customReader("data/count_100.csv"),
	)
}

func customReader(file string) *multicsv.LazyReader {
	return &multicsv.LazyReader{
		Init: func() (*csv.Reader, error) {
			f, err := os.Open(file)
			if err != nil {
				return nil, err
			}

			// customize csv.Reader
			r := csv.NewReader(f)
			r.LazyQuotes = true

			return r, nil
		},
	}
}
```

## License

http://www.opensource.org/licenses/mit-license.php
