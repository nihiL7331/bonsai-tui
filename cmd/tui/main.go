package main

import (
	"bonsai-tui/internal/tui"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	appModel := tui.New()

	if _, err := tea.NewProgram(appModel).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
