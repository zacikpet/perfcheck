package perfcheck

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/zacikpet/perfcheck/pkg/parsers"
	"github.com/zacikpet/perfcheck/templates"
)

func generateBenchmark(out string, model parsers.Api) *os.File {

	tmpl := template.New("default")

	tmpl.Parse(templates.DefaultTemplate)

	err := os.MkdirAll(filepath.Dir(out), os.ModePerm)
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
		fmt.Fprintf(os.Stderr, "Failed to execute template %s\n", tmpl.Name())
		panic(err)
	}

	return file
}
