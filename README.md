# multicsv

[![build](https://github.com/zero-pkg/multicsv/actions/workflows/ci.yml/badge.svg)](https://github.com/zero-pkg/multicsv/actions/workflows/ci.yml)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/zero-pkg/multicsv/blob/master/LICENSE)

MultiCSV is a multi-reader that provides the logical concatenation of the input CSV readers.
The readers are read sequentially. Once all inputs have returned `EOF`, `Read` returns `EOF`.
If any reader returns a non-nil, non-EOF error, `Read` returns that error.

## Installation

`go get -u github.com/zero-pkg/multicsv`

## Usage

```go
r := multicsv.NewReader(
    multicsv.LazyFileReader("data/users.csv"),
    multicsv.LazyFileReader("data/users2.csv", multicsv.WithSkipHeader()),
)

records, err := r.ReadAll()
if err != nil {
    panic(err)
}

fmt.Println(records)
```

## Extending `LazyReader`

```go
func main() {
	r := multicsv.NewReader(
		customReader("data/count_10.csv"),
		customReader("data/count_100.csv"),
	)
	_ = r
}

func customReader(file string) *multicsv.LazyReader {
	return &multicsv.LazyReader{
		Init: func() (*csv.Reader, error) {
			f, err := os.Open(file)
			if err != nil {
				return nil, err
			}

			// Customize csv.Reader.
			r := csv.NewReader(f)
			r.LazyQuotes = true

			return r, nil
		},
	}
}
```

## License

[MIT](https://opensource.org/license/mit)
