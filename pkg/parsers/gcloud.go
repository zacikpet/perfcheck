package parsers

import (
	"context"
	"errors"
	"fmt"
	"os"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"google.golang.org/api/iterator"
)

func ParseGCloudSLOs(projectId string, serviceId string, serviceUrl string, docsUrl string) Api {
	ctx := context.Background()

	monitoringClient, err := monitoring.NewServiceMonitoringClient(ctx)
	check(err)

	monitoringReq := &monitoringpb.ListServiceLevelObjectivesRequest{
		Parent: fmt.Sprintf("projects/%s/services/%s", projectId, serviceId),
	}

	slos := monitoringClient.ListServiceLevelObjectives(ctx, monitoringReq)

	var latencies []Metric
	var errorRates []Metric

	for {
		slo, err := slos.Next()
		if err == iterator.Done {
			break
		}
		check(err)

		latency, errorRate, err := parseSLO(*slo)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if latency != nil {
			latencies = append(latencies, *latency)
		}

		if errorRate != nil {
			errorRates = append(errorRates, *errorRate)
		}
	}

	return Api{
		BaseUrl: serviceUrl,
		Paths: []Path{
			{
				Method:   "GET",
				Pathname: "/",
				Detail: PathDetail{
					Latency:   latencies,
					ErrorRate: errorRates,
				},
			},
		},
	}
}

func parseSLO(slo monitoringpb.ServiceLevelObjective) (*Metric, *Metric, error) {

	sli := slo.ServiceLevelIndicator.Type

	basicSli, ok := sli.(*monitoringpb.ServiceLevelIndicator_BasicSli)
	if !ok {
		return nil, nil, errors.New("SLI is not basic_sli")
	}

	latency := basicSli.BasicSli.GetLatency()
	errorRate := basicSli.BasicSli.GetAvailability()

	if latency != nil {
		metric := Metric(fmt.Sprintf("p(%f) < %d", slo.Goal, latency.Threshold.AsDuration().Milliseconds()))
		return &metric, nil, nil
	} else if errorRate != nil {
		metric := Metric(fmt.Sprintf("rate < %f", 1-slo.Goal))
		return nil, &metric, nil
	}

	return nil, nil, errors.New("unknown SLI")
}
