package perfcheck

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/zacikpet/perfcheck/pkg/parsers"
)

func generateBenchmark(in string, out string, model parsers.Api) *os.File {

	tmpl, err := template.ParseFiles(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Template %s not found.\n", in)
		panic(err)
	}

	err = os.MkdirAll(filepath.Dir(out), os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory %s.\n", filepath.Dir(out))
		panic(err)
	}

	file, err := os.Create(out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create benchmark file %s\n", out)
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, model)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute benchmark %s\n", in)
	}

	return file
}
