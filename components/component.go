package components

import tea "github.com/charmbracelet/bubbletea"

type Component interface {
	Update(msg tea.Msg) (Component, tea.Cmd)
	View() string
	Key() string
	Value() interface{}
	Focus() tea.Cmd
	Unfocus() tea.Cmd
}
