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
defer r.Close()

records, err := r.ReadAll()
if err != nil {
    panic(err)
}

fmt.Println(records)
```

You can also read CSV data from bytes with the same options:

```go
r := multicsv.NewReader(
    multicsv.BytesReader(data, multicsv.WithSkipHeader()),
)
```

Use `WithAutoDetectDelimiter` when a file can use a common delimiter such as `,`, `;`, tab, or `|`:

```go
r := multicsv.NewReader(
    multicsv.LazyFileReader("data/users.csv", multicsv.WithAutoDetectDelimiter()),
)
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

func customReader(path string) *multicsv.LazyReader {
	var file *os.File

	return &multicsv.LazyReader{
		Init: func() (*csv.Reader, error) {
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}

			// Customize csv.Reader.
			r := csv.NewReader(f)
			r.LazyQuotes = true
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
```

## License

[MIT](https://opensource.org/license/mit)
