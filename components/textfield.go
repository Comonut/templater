package components

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3399ff"))
var blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

type TextField struct {
	Param string
	Title string
	Input textinput.Model
}

func (t TextField) Focus() tea.Cmd {
	panic("implement me")
}

func (t TextField) Unfocus() tea.Cmd {
	panic("implement me")
}

func (t TextField) Update(msg tea.Msg) (Component, tea.Cmd) {
	i, command := t.Input.Update(msg)
	t.Input = i
	return t, command
}

func (t TextField) View() string {
	return fmt.Sprintf("%s\n%s", t.Title, t.Input.View())
}

func (t TextField) Key() string {
	return t.Param
}

func (t TextField) Value() interface{} {
	return t.Input.Value()
}

func NewTextField(param, title, placeholder string) *TextField {
	t := textinput.NewModel()
	t.Placeholder = placeholder
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle

	return &TextField{
		Param: param,
		Title: title,
		Input: t,
	}
}
