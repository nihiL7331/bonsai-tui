package tui

import (
	"bonsai-tui/internal/tui/tabs"
	"bonsai-tui/internal/tui/tabs/build"
	"bonsai-tui/internal/tui/tabs/log"

	tea "charm.land/bubbletea/v2"
)

type AppModel struct {
	activeTab int
	tabs      []tabs.Tab
}

func New() AppModel {
	return AppModel{
		activeTab: 0,
		tabs: []tabs.Tab{
			build.New(),
			log.New(),
		},
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.activeTab++
			m.activeTab %= len(m.tabs)
			return m, nil
		case "shift+tab":
			m.activeTab--
			m.activeTab %= len(m.tabs)
			return m, nil
		}

		updatedModel, cmd := m.tabs[m.activeTab].Update(msg)
		m.tabs[m.activeTab] = updatedModel.(tabs.Tab)

		return m, cmd
	}

	for i, tab := range m.tabs {
		updatedTab, cmd := tab.Update(msg)
		m.tabs[i] = updatedTab.(tabs.Tab)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m AppModel) View() tea.View {
	var header string
	for i, tab := range m.tabs {
		if i == m.activeTab {
			header += "[ *" + tab.Title() + "* ] "
		} else {
			header += "[ " + tab.Title() + " ] "
		}
	}

	body := m.tabs[m.activeTab].View()
	body.SetContent(header + "\n\n" + body.Content)

	return body
}
