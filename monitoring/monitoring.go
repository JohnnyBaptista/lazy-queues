package monitoring

import (
	"fmt"
	"lazy-queues/client"
)

func FetchMetricDescriptors() {
	response, err := client.FetchData[any]("projects/" + client.ProjectName + "/metricDescriptors?filter=metric.type=starts_with(\"pubsub.googleapis.com\")" )

	if err != nil {
		fmt.Println("Error fetching metricDescriptors", err)
	}
	fmt.Println(response)

}
