package monitoring

import (
	"fmt"
	"lazy-queues/client"
	"lazy-queues/util"
	"net/url"
	"time"
)

const pubsubMetricBase = "pubsub.googleapis.com"

func FetchMetricDescriptors() {
	response, err := client.FetchData[any]("projects/" + client.ProjectName + "/metricDescriptors?filter=metric.type=starts_with(\"pubsub.googleapis.com\")")
	if err != nil {
		fmt.Println("Error fetching metricDescriptors", err)
	}
	fmt.Println(response)
}

func FetchTimeSeries[T any](metric metricType, start time.Time, end time.Time, subscriptionID string) (T, error) {
	metricPath := ResolveMetricPath(metric)
	util.Log.Info(
		"FetchTimeSeries",
		"metric", metricPath,
		"subscriptionID",
		subscriptionID,
		"start",
		start.String(),
		"end",
		end.String(),
	)

	isoStart := start.Format(time.RFC3339)
	isoEnd := end.Format(time.RFC3339)

	params := url.Values{}

	params.Set("filter",
		fmt.Sprintf(
			`metric.type="%s/%s" AND resource.labels.subscription_id="%s"`,
			pubsubMetricBase,
			metricPath,
			subscriptionID,
		))
	params.Set("interval.endTime", isoStart)
	params.Set("interval.startTime", isoEnd)

	path := "projects/" + client.ProjectName + "/timeSeries?" + params.Encode()

	util.Log.Info("FetchTimeSeries final path", "path", path)

	response, err := client.FetchData[T](path)
	if err != nil {
		var zero T
		return zero, err
	}

	return response, nil
}

func FetchOldestUnackedMessages() {
	currentMetric := SubscriptionMetricOldestUnackedMessageAge
	currentSubscription := "mercado-libre-answer-question-dlq-sub"

	now := time.Now()
	lastWeek := now.AddDate(0, 0, -7)

	response, err := FetchTimeSeries[any](currentMetric, now, lastWeek, currentSubscription)
	if err != nil {
		fmt.Println("Error fetching oldest unacked messages", err)
	}
	fmt.Println(response)
}
