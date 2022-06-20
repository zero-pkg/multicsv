package main

import (
	"fmt"

	"github.com/zero-pkg/multicsv"
)

func main() {
	r := multicsv.NewReader(
		multicsv.LazyFileReader("data/users.csv", true),
		multicsv.LazyFileReader("data/users2.csv", true),
	)

	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	fmt.Println(records)
}
