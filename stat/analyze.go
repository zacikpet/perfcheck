package stat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aclements/go-moremath/stats"
	"github.com/zacikpet/perf-check/parsers"
)

const alpha = 0.05

type DataPoint struct {
	Metric string `json:"metric"`
	Type   string `json:"type"`
	Data   Data   `json:"data"`
}

type Data struct {
	Time  string   `json:"time"`
	Value float64  `json:"value"`
	Tags  DataTags `json:"tags"`
}

type DataTags struct {
	ExpectedResponse string  `json:"expected_response"`
	Group            *string `json:"group"`
}

func AnalyzeData(dataFile string, model parsers.Api) {
	file, err := os.Open(dataFile)
	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	var points []DataPoint

	for scanner.Scan() {
		var point DataPoint
		json.Unmarshal(scanner.Bytes(), &point)

		points = append(points, point)
	}

	AnalyzeMetric(points, "http_req_duration", "::/")
}

func AnalyzeMetric(points []DataPoint, metric string, group string) {
	var data []float64

	for _, point := range points {
		if point.Metric == metric && point.Data.Tags.Group != nil && *point.Data.Tags.Group == group {
			data = append(data, point.Data.Value)
		}
	}

	res, err := stats.OneSampleTTest(makeSample(data), 100, stats.LocationLess)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.P)

	if res.P < alpha {
		fmt.Println("Null hypothesis rejected. CI/CD pass.")
	} else {
		fmt.Println("Null hypothesis not rejected. CI/CD fail.")
	}

}
