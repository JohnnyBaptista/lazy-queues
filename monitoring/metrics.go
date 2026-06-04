// Package monitoring lista as métricas padrão do PubSub.
// Segue o padrão pubsub.googleapis.com/METRIC_PATH.
// Documentação: https://cloud.google.com/monitoring/api/metrics_gcp
package monitoring

type metricType string

const (
	SubscriptionMetricAckMessageCount         metricType = "ack_message_count"
	SubscriptionMetricDLQMessageCount         metricType = "dead_letter_message_count"
	SubscriptionMetricOldestUnackedMessageAge metricType = "oldest_unacked_message_age"
)

func ResolveMetricPath(metric metricType) string {
	path := "subscription/" + string(metric)
	return path
}
