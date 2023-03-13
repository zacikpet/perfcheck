package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/pb33f/libopenapi"

	"github.com/zacikpet/perf-check/parsers"
	"github.com/zacikpet/perf-check/stat"
)

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {

	// stat.AnalyzeData("test.jsonl")

	// return
	godotenv.Load(".env")

	docsUrl := os.Getenv("DOCS_URL")

	res, err := http.Get(docsUrl)
	check(err)

	body, err := io.ReadAll(res.Body)
	check(err)

	document, err := libopenapi.NewDocument(body)
	check(err)

	model := parsers.ParseOpenAPI(document)

	tmpl, err := template.ParseFiles("templates/benchmark.js.tmpl")
	check(err)

	err = os.MkdirAll("benchmarks", os.ModePerm)
	check(err)

	file, err := os.Create("benchmarks/benchmark.js")
	check(err)
	defer file.Close()

	err = tmpl.Execute(file, model)
	check(err)

	fmt.Println("Benchmark generated.")

	_, err = exec.LookPath("k6")
	check(err)

	dataFile := "test.jsonl"

	cmd := exec.Command("k6", "run", file.Name(), "--out", fmt.Sprintf("json=%s", dataFile))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	stat.AnalyzeData("test.jsonl", model)

	if err != nil {
		fmt.Println("k6 threshold did not pass")
	} else {
		fmt.Println("k6 fine")
	}

}
