package tabs

import tea "charm.land/bubbletea/v2"

type Tab interface {
	tea.Model
	Title() string
}
