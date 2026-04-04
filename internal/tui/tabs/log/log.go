package log

import (
	"bonsai-tui/internal/tui/tabs"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
}

func New() tabs.Tab {
	return Model{}
}

func (m Model) Title() string { return "Log" }

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m Model) View() tea.View { return tea.NewView("Log Tab") }
