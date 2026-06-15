package monitoring

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"lazy-queues/client"
	"lazy-queues/state"
	"lazy-queues/util"
)

const pubsubMetricBase = "pubsub.googleapis.com"

func fetchTimeSeries[T TimeSeriesReponse](metric metricType, start time.Time, end time.Time, subscriptionID string) (T, error) {
	metricPath := ResolveMetricPath(metric)
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

	path := client.GoogleMonitoringBaseURL + "projects/" + state.State.ProjectID + "/timeSeries?" + params.Encode()

	response, err := client.FetchData[T](path)
	if err != nil {
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

		var zero T
		return zero, err
	}

	return response, nil
}

// FetchGenericMetric - pode ser 1 dia, 3 dias,  7 dias, 20 dias
func FetchGenericMetric(metric metricType, subscriptionID string, period int) ([]TimeSeries, error) {
	end := time.Now()
	start := end.AddDate(0, 0, -period)

	response, err := fetchTimeSeries(
		metric,
		start,
		end,
		subscriptionID,
	)
	if err != nil {
		var zero []TimeSeries
		util.Log.Error("Error fetching oldest unacked messages", "err", err)
		return zero, err
	}

	return response.TimeSeries, nil
}

func FetchMetrics(metrics []metricType, subscriptionID string, period int) map[metricType][]TimeSeries {
	timeSeriesByMetric := make(map[metricType][]TimeSeries)

	for _, met := range metrics {
		response, err := FetchGenericMetric(met, subscriptionID, period)
		if err != nil {
			util.Log.Error("Error fetching oldest unacked messages", "err", err, "metric", string(met))
		} else {
			timeSeriesByMetric[met] = response
		}
	}
	return timeSeriesByMetric
}

type Subscriptions struct {
	Name  string `json:"name"`
	Topic string `json:"topic"`
}

type SubscriptionList struct {
	Subscriptions []Subscriptions `json:"subscriptions"`
}

func GetSubscriptions() ([]string, error) {
	projectID := state.State.ProjectID
	path := client.GooglePubSubBaseURL + "projects/" + projectID + "/subscriptions"

	response, err := client.FetchData[SubscriptionList](path)
	if err != nil {
		return nil, err
	}

	var subscriptionList []string
	for _, subs := range response.Subscriptions {
		splitted := strings.Split(subs.Name, "/")
		subscriptionName := splitted[len(splitted)-1]
		subscriptionList = append(subscriptionList, subscriptionName)
	}

	return subscriptionList, nil
}
