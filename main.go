package main

import (
	"os"

	"lazy-queues/client"
	"lazy-queues/monitoring"
	"lazy-queues/state"
	"lazy-queues/tui"
	"lazy-queues/util"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	appState := &state.State
	theme := huh.ThemeBase()
	theme.Focused.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffe0000"))
	theme.Focused.Base = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("67")).
		Padding(0, 1)

	client.Execute()

	tui.StartForm(&appState.SubscriptionID, &appState.ProjectID, theme)

	if valid := appState.Validate(); !valid {
		os.Exit(1)
	}

	// timeSeriesResponse := monitoring.FetchMetrics(
	// 	monitoring.MetricsList[:],
	// 	state.State.SubscriptionID,
	// 	7,
	// )
	//
	// firstTimeSeries := timeSeriesResponse[monitoring.SubscriptionMetricOldestUnackedMessageAge]
	// tui.StartView(firstTimeSeries[0])
	timeSeriesResponse, err := monitoring.FetchGenericMetric(
		monitoring.SubscriptionMetricAckMessageCount,
		appState.SubscriptionID,
		7,
	)
	if err != nil {
		util.Log.Error("Erro buscando metrica")
		os.Exit(1)
	}

	tui.StartView(timeSeriesResponse[0])
}
