package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Api struct {
	BaseUrl string
	Paths   []Path
	Config  Config
}

type Config struct {
	Users    *int
	Duration *int
	Stages   []Stage
}

type Stage struct {
	Duration string
	Target   int
}

type Path struct {
	Method   string
	Pathname string
	Detail   PathDetail
}

type PathDetail struct {
	Latency      []Metric
	ErrorRate    []Metric
	ResponseSize []Metric
	Params       Params
}

type Metric string

type Examples struct {
	Path *map[string]any
}

type Params struct {
	Path  map[string]ParamDescription
	Query map[string]ParamDescription
}

type ParamDescription struct {
	Examples []any
	Pattern  *string
	Range    *RangeDescription
}

type RangeDescription struct {
	Min  int
	Max  int
	Step int
}

func (metric Metric) IsK6Supported() bool {
	return !strings.Contains(string(metric), "avg_stat")
}

func parseJSON[T any](dict any, target *T) {
	js, err := json.Marshal(dict)
	check(err)

	dec := json.NewDecoder(bytes.NewReader(js))
	dec.DisallowUnknownFields()

	err = dec.Decode(&target)
	check(err)
}

func check(errs ...error) {

	if len(errs) > 0 && errs[0] != nil {

		for i := range errs {
			fmt.Println(errs[i])
		}

		panic(errs[0])
	}
}
