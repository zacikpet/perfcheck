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
)

func main() {

	godotenv.Load(".env")

	docsUrl := os.Getenv("DOCS_URL")

	res, err := http.Get(docsUrl)
	check(err)

	body, err := io.ReadAll(res.Body)
	check(err)

	document, err := libopenapi.NewDocument(body)
	check(err)

	model := ParseOpenAPI(document)

	tmpl, err := template.ParseFiles("templates/benchmark.js.tmpl")
	check(err)

	err = os.MkdirAll("benchmarks", os.ModePerm)
	check(err)

	file, err := os.Create("benchmarks/benchmark.js")
	check(err)
	defer file.Close()

	tmpl.Execute(file, model)

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
