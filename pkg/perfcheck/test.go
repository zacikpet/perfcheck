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

	benchmark := generateBenchmark(outFile, model)

	if !noK6 {
		ok := RunK6(benchmark, k6DataFile)
		if !ok {
			return errors.New("k6: service does not conform to the service-level objectives")
		}
	}

	statOk := stat.AnalyzeData(outFile, model)

	if !statOk {
		return errors.New("perfcheck/stat: service does not conform to the service-level objectives")
	}

	return nil
}
