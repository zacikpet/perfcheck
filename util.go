package main

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-json"
)

func check(errs ...error) {
	if len(errs) > 0 && errs[0] != nil {

		for i := range errs {
			fmt.Println(errs[i])
		}

		panic(errs[0])
	}
}

func parseJSON[T any](dict any, target *T) {
	js, err := json.Marshal(dict)
	check(err)

	dec := json.NewDecoder(bytes.NewReader(js))
	dec.DisallowUnknownFields()

	err = dec.Decode(&target)
	check(err)
}
