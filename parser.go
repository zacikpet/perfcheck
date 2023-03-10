package main

import (
	"strings"

	"github.com/pb33f/libopenapi"
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
	SLO      SLO
}

type SLO struct {
	Latency   *int
	ErrorRate *float64
	Params    Params
}

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

func ParseOpenAPI(document libopenapi.Document) Api {
	version := document.GetVersion()

	if strings.HasPrefix(version, "2.") {
		return parseOpenAPIv2(document)
	}

	panic("Unsupported OpenAPI version")
}

func parseOpenAPIv2(document libopenapi.Document) Api {
	model, errs := document.BuildV2Model()
	check(errs...)

	var paths []Path

	for pathname := range model.Model.Paths.PathItems {
		operations := model.Model.Paths.PathItems[pathname].GetOperations()

		for method, operation := range operations {

			dict := operation.Extensions["x-perf-check"]

			var slo SLO
			parseJSON(dict, &slo)

			paths = append(paths, Path{
				Method:   method,
				Pathname: pathname,
				SLO:      slo,
			})
		}
	}

	_config := model.Model.Extensions["x-perf-check"]

	var config Config

	parseJSON(_config, &config)

	schemes := model.Model.Schemes

	if len(schemes) == 0 {
		panic("You must include at least one scheme (http, https)")
	}

	return Api{Paths: paths, BaseUrl: schemes[0] + "://" + model.Model.Host, Config: config}
}
