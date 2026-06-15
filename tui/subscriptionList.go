package tui

import (
	"fmt"
	"os"

	"lazy-queues/monitoring"

	tea "github.com/charmbracelet/bubbletea"
)

type subscriptionListViewModel struct {
	subscriptionOptions []string
	cursor              int
	selected            string
}

func subscriptionModel(subscriptionList []string) subscriptionListViewModel {
	return subscriptionListViewModel{
		subscriptionOptions: subscriptionList,
		cursor:              0,
		selected:            "",
	}
}

func (m subscriptionListViewModel) View() string {
	message := "Select what subscription to view: \n"

	for i, opt := range m.subscriptionOptions {
		curCursor := " "
		if m.cursor == i {
			curCursor = ">" // cursor!
		}

		message += fmt.Sprintf("%s %s\n", curCursor, opt)
	}

	message += "\nPress q to quit.\n"

	return message
}

func (m subscriptionListViewModel) Init() tea.Cmd {
	return tea.ClearScreen
}

func ListSubscriptionsView(options []string) {
	p := tea.NewProgram(subscriptionModel(options))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func chosenView(m subscriptionListViewModel) {
	// fetch timeSeries
	timeSeriesResponse := monitoring.FetchMetrics(
		monitoring.MetricsList[:],
		m.selected,
		7,
	)
	firstTimeSeries := timeSeriesResponse[monitoring.SubscriptionMetricOldestUnackedMessageAge]
	StartView(firstTimeSeries[0])
}

func (m subscriptionListViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.subscriptionOptions)-1 {
				m.cursor++
			}
		case "enter", "space":
			m.selected = m.subscriptionOptions[m.cursor]
			chosenView(m)
		}
	}
	return m, nil
}
