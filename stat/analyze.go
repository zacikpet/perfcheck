package stat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/aclements/go-moremath/stats"
	"github.com/zacikpet/perfcheck/parsers"
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

	for _, path := range model.Paths {

		for _, metric := range path.Detail.Latency {
			AnalyzeMetric(metric, points, path.Pathname, "http_req_duration")
		}

		for _, metric := range path.Detail.ErrorRate {
			AnalyzeMetric(metric, points, path.Pathname, "http_req_failed")
		}
	}
}

func AnalyzeMetric(metric parsers.Metric, points []DataPoint, pathname string, metricName string) {
	isAvgStat := regexp.MustCompile(`^\s*avg_stat\s*<\s*(\d+\.?\d*)\s*$`)

	submatches := isAvgStat.FindStringSubmatch(string(metric))

	if len(submatches) == 0 {
		fmt.Printf("no match\n")
	} else {
		avgStat, err := strconv.ParseFloat(submatches[1], 64)
		if err != nil {
			panic(err)
		}

		fmt.Printf("match %f\n", avgStat)
		AnalyzeAvgStat(points, metricName, fmt.Sprintf("::%s", pathname))
	}
}

func AnalyzeAvgStat(points []DataPoint, name string, group string) {
	var data []float64

	for _, point := range points {
		if point.Metric == name && point.Data.Tags.Group != nil && *point.Data.Tags.Group == group {
			data = append(data, point.Data.Value)
		}
	}

	res, err := stats.OneSampleTTest(makeSample(data), 100, stats.LocationLess)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to perform one sample t-test. Reason: %s\n", err)
		return
	}

	if res.P < alpha {
		fmt.Printf("%s{%s} Null hypothesis rejected. CI/CD pass. p = %2f\n", name, group, res.P)
	} else {
		fmt.Printf("%s{%s} Null hypothesis not rejected. CI/CD fail. p = %2f\n", name, group, res.P)
	}

}
