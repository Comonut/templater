package components

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3399ff"))
var blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
var noStyle = lipgloss.NewStyle()

type TextField struct {
	Param   string
	Title   string
	Input   textinput.Model
	focused bool
}

func (t TextField) Focus() (Component, tea.Cmd) {
	focus := t.Input.Focus()
	t.focused = true
	t.Input.PromptStyle = focusedStyle
	t.Input.TextStyle = focusedStyle
	return t, focus
}

func (t TextField) Unfocus() (Component, tea.Cmd) {
	t.Input.Blur()
	t.focused = false
	t.Input.PromptStyle = noStyle
	t.Input.TextStyle = noStyle
	return t, nil
}

func (t TextField) Update(msg tea.Msg) (Component, tea.Cmd) {
	i, command := t.Input.Update(msg)
	t.Input = i
	return t, command
}

func (t TextField) View() string {
	var title string
	if t.focused {
		title = focusedStyle.Copy().Render(t.Title)
	} else {
		title = t.Title
	}
	return fmt.Sprintf("%s\n%s", title, t.Input.View())
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
