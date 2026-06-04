package monitoring

import (
	"fmt"
	"lazy-queues/client"
	"net/url"
	"time"
)

const pubsubMetricBase = "pubsub.googleapis.com/subscription"

func FetchMetricDescriptors() {
	response, err := client.FetchData[any]("projects/" + client.ProjectName + "/metricDescriptors?filter=metric.type=starts_with(\"pubsub.googleapis.com\")")
	if err != nil {
		fmt.Println("Error fetching metricDescriptors", err)
	}
	fmt.Println(response)
}

func FetchOldestUnackedMessages() {
	currentMetric := SubscriptionMetricOldestUnackedMessageAge

	now := time.Now()
	lastWeek := now.AddDate(0, 0, -7)
	isoNow := now.Format(time.RFC3339)
	isoLastWeek := lastWeek.Format(time.RFC3339)

	params := url.Values{}

	params.Set("filter",
		fmt.Sprintf(`metric.type="%s/%s" AND resource.labels.subscription_id="mercado-libre-answer-question-analytics"`,
			pubsubMetricBase,
			currentMetric))
	params.Set("interval.endTime", isoNow)
	params.Set("interval.startTime", isoLastWeek)

	path := "projects/" + client.ProjectName + "/timeSeries?" + params.Encode()
	fmt.Println(path)

	response, err := client.FetchData[any](path)
	if err != nil {
		fmt.Println("Error fetching oldest unacked messages", err)
	}
	fmt.Println(response)
}
