package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zacikpet/perfcheck/pkg/perfcheck"
)

func main() {

	app := &cli.App{
		Name:  "perfcheck",
		Usage: "Automatic benchmarks of APIs",
		Commands: []*cli.Command{
			{
				Name:  "get-template",
				Usage: "Output default template to file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "outFile",
						Aliases: []string{"o"},
						Usage:   "Name of output file",
						Value:   "benchmark.js.tmpl",
						EnvVars: []string{"OUT_FILE"},
					},
				},
				Action: func(ctx *cli.Context) error {
					return perfcheck.GetTemplate(ctx.String("outFile"))
				},
			},
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "test a service",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "source",
						Aliases: []string{"s"},
						Usage:   "Source of the SLOs. Allowed: openapi, gcloud",
						EnvVars: []string{"SOURCE"},
					},
					&cli.StringFlag{
						Name:    "docsUrl",
						Aliases: []string{"d"},
						Usage:   "URL of the OpenAPI documentation",
						EnvVars: []string{"DOCS_URL"},
					},
					&cli.StringFlag{
						Name:    "gcloudProjectId",
						Usage:   "Google Cloud Project ID",
						EnvVars: []string{"GCLOUD_PROJECT_ID"},
					},
					&cli.StringFlag{
						Name:    "gcloudServiceId",
						Usage:   "Google Cloud Service ID",
						EnvVars: []string{"GCLOUD_SERVICE_ID"},
					},
					&cli.StringFlag{
						Name:    "gcloudServiceUrl",
						Usage:   "Google Cloud Service URL",
						EnvVars: []string{"GCLOUD_SERVICE_URL"},
					},
					&cli.StringFlag{
						Name:    "templateFile",
						Value:   "",
						Usage:   "Template file for the benchmark",
						EnvVars: []string{"TEMPLATE_FILE"},
					},
					&cli.StringFlag{
						Name:    "outFile",
						Value:   "benchmarks/benchmark.js",
						Usage:   "Output file for the k6 benchmark",
						EnvVars: []string{"OUT_FILE"},
					},
					&cli.StringFlag{
						Name:    "k6DataFile",
						Value:   "k6.jsonl",
						Usage:   "Output file for the k6 benchmark data",
						EnvVars: []string{"K6_DATA_FILE"},
					},
					&cli.BoolFlag{
						Name:    "no-k6",
						Value:   false,
						Usage:   "Don't run k6 automatically",
						EnvVars: []string{"NO_K6"},
					},
				},
				Action: func(ctx *cli.Context) error {
					return perfcheck.Test(
						ctx.String("source"),
						ctx.String("docsUrl"),
						ctx.String("gcloudProjectId"),
						ctx.String("gcloudServiceId"),
						ctx.String("gcloudServiceUrl"),
						ctx.String("templateFile"),
						ctx.String("outFile"),
						ctx.String("k6DataFile"),
						ctx.Bool("no-k6"),
					)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
