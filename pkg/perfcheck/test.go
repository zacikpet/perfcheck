package perfcheck

import (
	"fmt"

	"github.com/zacikpet/perfcheck/pkg/parsers"
	"github.com/zacikpet/perfcheck/pkg/stat"
)

func Test(
	source string,
	docsUrl string,
	projectId string,
	serviceId string,
	serviceUrl string,
	template string,
	outFile string,
	k6DataFile string,
) {
	var model parsers.Api

	switch source {
	case "openapi":
		document := parsers.FetchOpenAPI(docsUrl)
		model = parsers.ParseOpenAPI(document)

	case "gcloud":
		model = parsers.ParseGCloudSLOs(projectId, serviceId, serviceUrl, docsUrl)

	default:
		panic(fmt.Sprintf("Invalid source %s", source))
	}

	benchmark := generateBenchmark(template, outFile, model)

	k6Ok := RunK6(benchmark, k6DataFile)

	stat.AnalyzeData(outFile, model)

	if k6Ok {
		fmt.Println("k6 fine")
	} else {
		fmt.Println("k6 threshold did not pass")
	}
}
