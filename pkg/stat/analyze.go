package stat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/aclements/go-moremath/stats"
	"github.com/zacikpet/perfcheck/pkg/parsers"
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

func AnalyzeData(dataFile string, model parsers.Api) bool {
	file, err := os.Open(dataFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open data file %s\n", dataFile)
		panic(err)
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

	allOk := true

	for _, path := range model.Paths {

		for _, metric := range path.Detail.Latency {
			isStat, ok := AnalyzeMetric(metric, points, path.Pathname, "http_req_duration")
			if !isStat {
				continue
			}

			if ok {
				fmt.Printf("%s Latency(%s): OK\n", path.Pathname, metric)
			} else {
				allOk = false
				fmt.Printf("%s Latency(%s): Failed\n", path.Pathname, metric)
			}
		}

		for _, metric := range path.Detail.ErrorRate {
			isStat, ok := AnalyzeMetric(metric, points, path.Pathname, "http_req_failed")
			if !isStat {
				continue
			}
			if ok {
				fmt.Printf("%s Error rate(%s): OK\n", path.Pathname, metric)
			} else {
				allOk = false
				fmt.Printf("%s Error rate(%s): Failed\n", path.Pathname, metric)
			}
		}

		for _, metric := range path.Detail.ResponseSize {
			isStat, ok := AnalyzeMetric(metric, points, path.Pathname, "response_bytes")
			if !isStat {
				continue
			}
			if ok {
				fmt.Printf("%s Response size(%s): OK\n", path.Pathname, metric)
			} else {
				allOk = false
				fmt.Printf("%s Response size(%s): Failed\n", path.Pathname, metric)
			}
		}
	}

	return allOk
}

func AnalyzeMetric(metric parsers.Metric, points []DataPoint, pathname string, metricName string) (bool, bool) {
	isAvgStat := regexp.MustCompile(`^\s*avg_stat\s*<\s*(\d+\.?\d*)\s*$`)

	submatches := isAvgStat.FindStringSubmatch(string(metric))

	if len(submatches) == 0 {
		return false, true
	} else {
		avgStat, err := strconv.ParseFloat(submatches[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid metric %s\n", metric)
		}

		ok := AnalyzeAvgStat(points, metricName, fmt.Sprintf("::%s", pathname), avgStat)
		return true, ok
	}
}

func AnalyzeAvgStat(points []DataPoint, name string, group string, target float64) bool {
	var data []float64

	for _, point := range points {

		if point.Metric == name && point.Data.Tags.Group != nil && *point.Data.Tags.Group == group {
			data = append(data, point.Data.Value)
		}
	}

	res, err := stats.OneSampleTTest(makeSample(data), target, stats.LocationLess)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to perform one sample t-test. Reason: %s\n", err)
		return true
	}

	return res.P < alpha
}
