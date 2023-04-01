package perfcheck

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/zacikpet/perfcheck/pkg/parsers"
	"github.com/zacikpet/perfcheck/templates"
)

func GetTemplate(outFile string) error {
	err := os.MkdirAll(filepath.Dir(outFile), os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory %s.\n", filepath.Dir(outFile))
		return err
	}

	file, err := os.Create(outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create template file %s\n", outFile)
		return err
	}

	defer file.Close()

	_, err = file.WriteString(templates.DefaultTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write template to file %s\n", outFile)
		return err
	}

	return nil
}

func generateBenchmark(templateFile string, out string, model parsers.Api) *os.File {

	tmpl := getTemplate(templateFile)

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

func getDefaultTemplate() *template.Template {
	tmpl := template.New("default")

	tmpl.Parse(templates.DefaultTemplate)

	return tmpl
}

func getTemplate(templateFile string) *template.Template {
	if templateFile == "" {
		return getDefaultTemplate()
	} else {
		tmpl := template.New(templateFile)

		tmpl.Parse(templates.DefaultTemplate)
		return tmpl
	}
}
