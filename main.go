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

	tui.StartForm(&appState.ProjectID, theme)

	if valid := appState.Validate(); !valid {
		os.Exit(1)
	}

	results, err := monitoring.GetSubscriptions()
	if err != nil {
		util.Log.Error("Ops", "err", err)
	}

	tui.ListSubscriptionsView(results[0:25])
}
