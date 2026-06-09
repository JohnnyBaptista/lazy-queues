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

func FetchTimeSeries[T TimeSeriesReponse](metric metricType, start time.Time, end time.Time, subscriptionID string) (T, error) {
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
	params.Set("interval.startTime", isoStart)
	params.Set("interval.endTime", isoEnd)

	path := "projects/" + client.ProjectName + "/timeSeries?" + params.Encode()

	util.Log.Info("FetchTimeSeries final path", "path", path)

	response, err := client.FetchData[T](path)
	if err != nil {
		var zero T
		return zero, err
	}

	return response, nil
}

// FetchGenericMetric - pode ser 1 dia, 3 dias,  7 dias, 20 dias
func FetchGenericMetric(metric metricType, subscriptionID string, period int) {
	end := time.Now()
	start := end.AddDate(0, 0, -period)

	response, err := FetchTimeSeries(
		metric,
		start,
		end,
		subscriptionID,
	)
	if err != nil {
		fmt.Println("Error fetching oldest unacked messages", err)
	}
	if len(response.TimeSeries) > 0 && len(response.TimeSeries[0].Points) > 0 {
		ts := response.TimeSeries[0]
		fmt.Printf("Metric: %+v\n", ts.Metric)
		fmt.Printf("First point: %+v\n", *&ts.Points[0].Interval.EndTime)
		fmt.Printf("First point value: %+v\n", *ts.Points[0].Value.Int64Value)
	}
}

func FetchOldestUnackedMessages() {
	currentMetric := SubscriptionMetricOldestUnackedMessageAge
	currentSubscription := "mercado-libre-answer-question-dlq-sub"

	FetchGenericMetric(currentMetric, currentSubscription, 1)
}
