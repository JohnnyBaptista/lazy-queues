package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

func StartForm(projectID *string, theme *huh.Theme) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project ID").
				Value(projectID),
		),
	).WithTheme(theme)

	if err := form.Run(); err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}
}
