package log

import (
	"bonsai-tui/internal/config"
	"bonsai-tui/internal/tui/tabs"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	config   config.Config
	messages []string
	height   int
}

func New(cfg config.Config) tabs.Tab {
	return Model{
		config:   cfg,
		messages: make([]string, 0),
	}
}

func (m Model) Title() string { return "Log" }

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.height = sizeMsg.Height - 2
	}
	timeStr := time.Now().Format("15:04:05")
	logLine := fmt.Sprintf("[%s] %T: %+v", timeStr, msg, msg)

	m.messages = append([]string{logLine}, m.messages...)

	if len(m.messages) > m.height {
		m.messages = m.messages[:m.height]
	}

	return m, nil
}

func (m Model) View() tea.View {
	return tea.NewView(strings.Join(m.messages, "\n"))
}
