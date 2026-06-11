package tui

import (
	"fmt"
	"os"

	"lazy-queues/monitoring"
	"lazy-queues/state"

	"charm.land/lipgloss/v2"
	tslc "github.com/NimbleMarkets/ntcharts/v2/linechart/timeserieslinechart"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone/v2"
)

// renderizar formulario (initialForm)
// listar formularios
type viewsModel struct {
	appState    *state.AppState
	chart       tslc.Model
	zoneManager *zone.Manager
	choices     []string
	cursor      int
	Choice      map[int]struct{}
	Quitting    bool
	timeSeries  monitoring.TimeSeries
}

func (m viewsModel) Init() tea.Cmd {
	return nil
}

const (
	defaultChartWidth  = 30
	defaultChartHeight = 12
)

func initialModel(timeSeries monitoring.TimeSeries) viewsModel {
	zoneManager := zone.New()

	chart := tslc.New(defaultChartWidth, defaultChartHeight)
	chart.SetZoneManager(zoneManager)
	mapChart(timeSeries.Points, &chart)
	chart.DrawBrailleAll()

	return viewsModel{
		appState:   &state.State,
		timeSeries: timeSeries,
		choices: []string{
			string(monitoring.SubscriptionMetricOldestUnackedMessageAge),
			string(monitoring.SubscriptionMetricAckMessageCount),
			string(monitoring.SubscriptionMetricDLQMessageCount),
		},
		Choice:      make(map[int]struct{}),
		Quitting:    false,
		chart:       chart,
		zoneManager: zoneManager,
		cursor:      0,
	}
}

func (m viewsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	//
	// 	// Cool, what was the actual key pressed?
	// 	switch msg.String() {
	//
	// 	// These keys should exit the program.
	// 	case "ctrl+c", "q":
	// 		return m, tea.Quit
	//
	// 	// The "up" and "k" keys move the cursor up
	// 	case "up", "k":
	// 		if m.cursor > 0 {
	// 			m.cursor--
	// 		}
	//
	// 	// The "down" and "j" keys move the cursor down
	// 	case "down", "j":
	// 		if m.cursor < len(m.choices)-1 {
	// 			m.cursor++
	// 		}
	//
	// 	// The "enter" key and the space bar toggle the selected state
	// 	// for the item that the cursor is pointing at.
	// 	case "enter", "space":
	// 		_, ok := m.Choice[m.cursor]
	// 		if ok {
	// 			delete(m.Choice, m.cursor)
	// 		} else {
	// 			m.Choice[m.cursor] = struct{}{}
	// 		}
	// 	}
	// }

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.

	return graphicsView(m, msg)
}

func mapChart(points []monitoring.Point, chart *tslc.Model) {
	for _, point := range points {
		date := point.Interval.EndTime
		chart.Push(tslc.TimePoint{date, float64(point.Value.Int64Value)})
	}
}

func graphicsView(m viewsModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.chart.Resize(msg.Width-2, msg.Height-2)
	}

	// forward Bubble Tea Msg to time series chart
	// and draw all data sets using braille runes
	m.chart, _ = m.chart.Update(msg)
	m.chart.DrawBrailleAll()
	return m, nil
}

func (m viewsModel) View() string {
	return m.zoneManager.Scan(
		lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")). // purple
			Render(m.chart.View()),
	)
}

func StartView(timeSeries monitoring.TimeSeries) {
	p := tea.NewProgram(initialModel(timeSeries))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
