package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/pb33f/libopenapi"
)

type path struct {
	Pathname  string
	Latency   int
	ErrorRate float64
}

type model struct {
	Paths []path
}

func check(errs ...error) {
	if len(errs) > 0 && errs[0] != nil {

		for i := range errs {
			fmt.Println(errs[i])
		}

		panic(errs[0])
	}
}

func buildModel(document libopenapi.Document) model {
	version := document.GetVersion()

	var pathnames []string
	var slos []any

	if version[0] == '2' {
		model, errs := document.BuildV2Model()
		check(errs...)

		for pathname := range model.Model.Paths.PathItems {
			slo := model.Model.Paths.PathItems[pathname].Get.Extensions["x-perf-check"]
			slos = append(slos, slo)
			pathnames = append(pathnames, pathname)
		}

	} else if version[0] == '3' {
		model, errs := document.BuildV3Model()
		check(errs...)

		for pathname := range model.Model.Paths.PathItems {
			slo := model.Model.Paths.PathItems[pathname].Get.Extensions["x-perf-check"]
			slos = append(slos, slo)
			pathnames = append(pathnames, pathname)
		}

	} else {
		panic(fmt.Sprintf("Unsupported document version (required 2 or 3, found %s)", version))
	}

	var paths []path

	for i, slo := range slos {
		if slo == nil {
			continue
		}

		_slo, ok := slo.(map[string]any)

		if !ok {
			panic("x-perf-check has invalid type")
		}

		latency := _slo["latency"]
		errorRate := _slo["errorRate"]

		l, ok := latency.(int)

		if !ok {
			panic("latency must be a number")
		}

		e, ok := errorRate.(float64)

		if !ok {
			panic("error rate must be a float64")
		}

		paths = append(paths, path{Pathname: pathnames[i], Latency: l, ErrorRate: e})
	}

	return model{Paths: paths}

}

func main() {
	godotenv.Load(".env")

	docsUrl := os.Getenv("DOCS_URL")

	res, err := http.Get(docsUrl)
	check(err)

	body, err := io.ReadAll(res.Body)
	check(err)

	document, err := libopenapi.NewDocument(body)
	check(err)

	model := buildModel(document)

	tmpl, err := template.ParseFiles("templates/benchmark.tmpl")
	check(err)

	err = os.MkdirAll("benchmarks", os.ModePerm)
	check(err)

	file, err := os.Create("benchmarks/benchmark.js")
	check(err)
	defer file.Close()

	vars := make(map[string]interface{})

	vars["BaseUrl"] = os.Getenv("BASE_URL")
	vars["Paths"] = model.Paths

	tmpl.Execute(file, vars)

	fmt.Println("Benchmark generated.")
}
