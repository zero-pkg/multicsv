package main

import (
	"fmt"

	"github.com/zero-pkg/multicsv"
)

func main() {
	r := multicsv.NewReader(
		multicsv.LazyFileReader("data/users.csv", multicsv.WithSkipHeader(), multicsv.WithAutoDetectDelimiter()),
		multicsv.LazyFileReader("data/users2.csv", multicsv.WithSkipHeader(), multicsv.WithAutoDetectDelimiter()),
	)

	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		fmt.Printf("%+v\n", record)
	}
}
