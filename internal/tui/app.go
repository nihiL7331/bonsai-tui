package tui

import (
	"bonsai-tui/internal/config"
	"bonsai-tui/internal/tui/tabs"
	"bonsai-tui/internal/tui/tabs/build"
	"bonsai-tui/internal/tui/tabs/log"

	tea "charm.land/bubbletea/v2"
)

type AppModel struct {
	config    config.Config
	activeTab uint
	tabs      []tabs.Tab
}

func New(cfg config.Config) AppModel {
	return AppModel{
		config:    cfg,
		activeTab: 0,
		tabs: []tabs.Tab{
			build.New(cfg),
			log.New(cfg),
		},
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	for i, tab := range m.tabs {
		if tab.Title() == "Log" {
			updatedLog, _ := tab.Update(msg)
			m.tabs[i] = updatedLog.(tabs.Tab)
			break
		}
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.activeTab++
			m.activeTab %= uint(len(m.tabs))
			return m, nil
		case "shift+tab":
			m.activeTab--
			m.activeTab %= uint(len(m.tabs))
			return m, nil
		}

		if m.tabs[m.activeTab].Title() != "Log" {
			updatedModel, cmd := m.tabs[m.activeTab].Update(msg)
			m.tabs[m.activeTab] = updatedModel.(tabs.Tab)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)
	}

	for i, tab := range m.tabs {
		if tab.Title() == "Log" {
			continue
		}
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
		if uint(i) == m.activeTab {
			header += "[ *" + tab.Title() + "* ] "
		} else {
			header += "[ " + tab.Title() + " ] "
		}
	}

	body := m.tabs[m.activeTab].View()
	body.SetContent(header + "\n\n" + body.Content)

	return body
}
