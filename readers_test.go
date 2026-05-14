package multicsv

import (
	"encoding/csv"
	"os"
	"testing"
)

func TestLazyReader(t *testing.T) {
	r := &LazyReader{
		Init: func() (*csv.Reader, error) {
			f, err := os.Open("testdata/custom.csv")
			ok(t, err)

			reader := csv.NewReader(f)
			reader.LazyQuotes = true

			_, err = reader.Read() // skip header
			ok(t, err)

			return reader, nil
		},
	}

	var cnt int

	for {
		fields, err := r.Read()
		if err != nil {
			break
		}

		cnt++

		equals(t, 2, len(fields))
	}

	equals(t, 10, cnt)
}

func TestLazyFileReader(t *testing.T) {
	r := LazyFileReader("testdata/basic.csv")

	var cnt int

	for {
		fields, err := r.Read()
		if err != nil {
			break
		}

		cnt++

		equals(t, 6, len(fields))
	}

	equals(t, 6, cnt)
}

func TestLazyFileReaderSkipHeader(t *testing.T) {
	r := LazyFileReader("testdata/basic.csv", WithSkipHeader())

	var cnt int

	for {
		fields, err := r.Read()
		if err != nil {
			break
		}

		cnt++

		equals(t, 6, len(fields))
	}

	equals(t, 5, cnt)
}

func TestLazyFileReaderAutoDetectDelimiter(t *testing.T) {
	r := LazyFileReader("testdata/semicolon.csv", WithAutoDetectDelimiter(), WithSkipHeader())

	records, err := NewReader(r).ReadAll()
	ok(t, err)

	equals(t, [][]string{
		{"Ann", "Moscow", "likes, commas"},
		{"Bob", "Berlin", "uses; semicolons"},
	}, records)
}

func TestLazyFileReaderAutoDetectDelimiterKeepsComma(t *testing.T) {
	r := LazyFileReader("testdata/basic.csv", WithAutoDetectDelimiter())

	var cnt int

	for {
		fields, err := r.Read()
		if err != nil {
			break
		}

		cnt++

		equals(t, 6, len(fields))
	}

	equals(t, 6, cnt)
}

func TestBytesReader(t *testing.T) {
	r := BytesReader([]byte("name,city\nAnn,Moscow\nBob,Berlin\n"))

	records, err := NewReader(r).ReadAll()
	ok(t, err)

	equals(t, [][]string{
		{"name", "city"},
		{"Ann", "Moscow"},
		{"Bob", "Berlin"},
	}, records)
}

func TestBytesReaderSkipHeader(t *testing.T) {
	r := BytesReader([]byte("name,city\nAnn,Moscow\nBob,Berlin\n"), WithSkipHeader())

	records, err := NewReader(r).ReadAll()
	ok(t, err)

	equals(t, [][]string{
		{"Ann", "Moscow"},
		{"Bob", "Berlin"},
	}, records)
}

func TestBytesReaderAutoDetectDelimiter(t *testing.T) {
	data := []byte(
		"name;city;note\n" +
			"Ann;Moscow;\"likes, commas\"\n" +
			"Bob;Berlin;\"uses; semicolons\"\n",
	)
	r := BytesReader(data, WithAutoDetectDelimiter(), WithSkipHeader())

	records, err := NewReader(r).ReadAll()
	ok(t, err)

	equals(t, [][]string{
		{"Ann", "Moscow", "likes, commas"},
		{"Bob", "Berlin", "uses; semicolons"},
	}, records)
}

func TestLazyFileReaderError(t *testing.T) {
	r := LazyFileReader("testdata/nonexists.csv")

	_, err := r.Read()
	assert(t, err != nil, "err is nil")
}
