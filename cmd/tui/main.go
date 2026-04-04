package main

import (
	"bonsai-tui/internal/config"
	"bonsai-tui/internal/tui"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to open bonsai.toml: %v\n", err)
		os.Exit(1)
	}

	appModel := tui.New(cfg)

	p := tea.NewProgram(appModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
