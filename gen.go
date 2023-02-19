package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/pb33f/libopenapi"
)

type SLO struct {
	Latency   int
	ErrorRate float64
}

type request struct {
	Method   string
	Pathname string
	SLO      SLO
}

type model struct {
	requests []request
}

func check(errs ...error) {
	if len(errs) > 0 && errs[0] != nil {

		for i := range errs {
			fmt.Println(errs[i])
		}

		panic(errs[0])
	}
}

func validateMetadata(metadata any) (*SLO, []error) {
	slo, ok := metadata.(map[string]any)

	var errs []error

	if !ok {
		errs = append(errs, errors.New("invalid type of x-perf-check"))
	}

	latency := slo["latency"]
	errorRate := slo["errorRate"]

	_latency, ok := latency.(int)

	if !ok {
		errs = append(errs, errors.New("latency must be a number"))
	}

	_errorRate, ok := errorRate.(float64)

	if !ok {
		errs = append(errs, errors.New("error rate must be a float64"))
	}

	return &SLO{Latency: _latency, ErrorRate: _errorRate}, errs
}

func buildModel(document libopenapi.Document) model {
	version := document.GetVersion()

	var methods []string
	var pathnames []string
	var slos []any

	if version[0] == '2' {
		model, errs := document.BuildV2Model()
		check(errs...)

		for pathname := range model.Model.Paths.PathItems {
			ops := model.Model.Paths.PathItems[pathname].GetOperations()

			for method, op := range ops {
				slos = append(slos, op.Extensions["x-perf-check"])
				pathnames = append(pathnames, pathname)
				methods = append(methods, method)
			}
		}

	} else if version[0] == '3' {
		model, errs := document.BuildV3Model()
		check(errs...)

		for pathname := range model.Model.Paths.PathItems {
			ops := model.Model.Paths.PathItems[pathname].GetOperations()

			for method, op := range ops {
				slos = append(slos, op.Extensions["x-perf-check"])
				pathnames = append(pathnames, pathname)
				methods = append(methods, method)
			}
		}

	} else {
		panic(fmt.Sprintf("Unsupported document version (required 2 or 3, found %s)", version))
	}

	var paths []request

	for i, slo := range slos {
		if slo == nil {
			continue
		}

		_slo, errs := validateMetadata(slo)
		check(errs...)

		paths = append(paths, request{
			Pathname: pathnames[i],
			Method:   methods[i],
			SLO:      *_slo,
		})
	}

	return model{requests: paths}
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

	tmpl, err := template.ParseFiles("templates/benchmark.js.tmpl")
	check(err)

	err = os.MkdirAll("benchmarks", os.ModePerm)
	check(err)

	file, err := os.Create("benchmarks/benchmark.js")
	check(err)
	defer file.Close()

	vars := make(map[string]interface{})

	vars["BaseUrl"] = os.Getenv("BASE_URL")
	vars["Paths"] = model.requests

	tmpl.Execute(file, vars)

	fmt.Println("Benchmark generated.")

	_, err = exec.LookPath("k6")
	check(err)

	cmd := exec.Command("k6", "run", file.Name())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		fmt.Println("✋")
	} else {
		fmt.Println("👌")
	}

}
