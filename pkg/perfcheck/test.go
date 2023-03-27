package perfcheck

import (
	"errors"
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
	noK6 bool,
) error {

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

	statOk := stat.AnalyzeData(outFile, model)

	if !noK6 && !k6Ok {
		return errors.New("k6: service does not conform to the service-level objectives")
	}

	if !statOk {
		return errors.New("perfcheck/stat: service does not conform to the service-level objectives")
	}

	return nil
}
