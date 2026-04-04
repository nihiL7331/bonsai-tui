package build

import (
	"bonsai-tui/internal/config"
	"bonsai-tui/internal/tui/tabs"
	"time"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	config     config.Config
	isBuilding bool
	logs       string
}

type buildFinishedMsg struct {
	success bool
	logs    string
}

func New(cfg config.Config) tabs.Tab {
	return Model{
		config:     cfg,
		isBuilding: false,
		logs:       "",
	}
}

func (m Model) Title() string { return "Build" }

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "b" && !m.isBuilding {
			m.isBuilding = true
			m.logs = "Building..."
			return m, runCompiler()
		}
	case buildFinishedMsg:
		m.isBuilding = false
		if msg.success {
			m.logs = msg.logs
		} else {
			m.logs = "Build failed!"
		}
		return m, nil
	}

	return m, nil
}

func (m Model) View() tea.View { return tea.NewView("Build Tab") }

func runCompiler() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)

		return buildFinishedMsg{
			success: true,
			logs:    "MSDF Font Generated. Atlas Packed.",
		}
	}
}
