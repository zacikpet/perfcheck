package parsers

import (
	"context"
	"fmt"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"google.golang.org/api/iterator"
)

func BuildGCloudURl(projectId string, serviceId string) string {
	return fmt.Sprintf(
		"https://monitoring.googleapis.com/v3/projects/%s/services/%s/serviceLevelObjectives",
		projectId,
		serviceId,
	)
}

func ParseGCloudSLOs(projectId string, serviceId string) Api {
	ctx := context.Background()

	client, err := monitoring.NewServiceMonitoringClient(ctx)
	check(err)

	req := &monitoringpb.ListServiceLevelObjectivesRequest{
		Parent: fmt.Sprintf("projects/%s/services/%s", projectId, serviceId),
	}

	slos := client.ListServiceLevelObjectives(ctx, req)

	for {
		slo, err := slos.Next()

		if err == iterator.Done {
			break
		}
		check(err)

		parseSLO(*slo)
	}

	return Api{}
}

func parseSLO(slo monitoringpb.ServiceLevelObjective) {
	fmt.Println(slo)
}
