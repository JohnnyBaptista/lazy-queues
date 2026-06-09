// Package monitoring lista as métricas padrão do PubSub.
// Segue o padrão pubsub.googleapis.com/METRIC_PATH.
// depois criar uma pasta para mapeamentos apenas do google
// Documentação: https://cloud.google.com/monitoring/api/metrics_gcp
package monitoring

import "time"

type metricType string

const (
	SubscriptionMetricAckMessageCount         metricType = "ack_message_count"
	SubscriptionMetricDLQMessageCount         metricType = "dead_letter_message_count"
	SubscriptionMetricOldestUnackedMessageAge metricType = "oldest_unacked_message_age" // Points.Value.Int64Value -> idade em segundos
)

func ResolveMetricPath(metric metricType) string {
	path := "subscription/" + string(metric)
	return path
}

// TimeSeries representa uma coleção de pontos de dados que descrevem os valores
// variantes no tempo de uma métrica. Identificada pela combinação de um recurso
// monitorado e uma métrica totalmente especificados.
// Ref: https://docs.cloud.google.com/monitoring/api/ref_v3/rest/v3/TimeSeries
type TimeSeries struct {
	Metric *Metric  `json:"metric"`
	Points []*Point `json:"points"`
}

type Metric struct {
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

type Point struct {
	Interval *TimeInterval `json:"interval"`
	Value    *TypedValue   `json:"value"`
}

type TimeInterval struct {
	EndTime time.Time `json:"endTime"`
}

type TypedValue struct {
	DoubleValue *float64 `json:"doubleValue,omitempty"`
	Int64Value  *int64   `json:"int64Value,omitempty,string"`
}

type TimeSeriesReponse struct {
	TimeSeries []*TimeSeries `json:"timeSeries"`
}
