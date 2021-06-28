package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type Choice struct {
	Param   string
	title   string
	options []string
	choice  int
}

func (c Choice) Update(msg tea.Msg) (Component, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if c.choice == -1 || c.choice == 0 {
				c.choice = len(c.options) - 1
				return c, nil
			}
			c.choice -= 1
			return c, nil
		case "right":
			if c.choice == -1 || c.choice == len(c.options)-1 {
				c.choice = 0
				return c, nil
			}
			c.choice += 1
			return c, nil
		}
	}
	return c, nil
}

func (c Choice) View() string {
	var builder strings.Builder
	builder.WriteString(c.title)
	builder.WriteString("\n")
	for _, option := range c.options {
		builder.WriteString("\n")
		builder.WriteString(option)
	}

	return builder.String()
}

func (c Choice) Key() string {
	return c.Param
}

func (c Choice) Value() interface{} {
	if c.choice == -1 {
		return ""
	}
	return c.options[c.choice]
}

func (c Choice) Focus() tea.Cmd {
	panic("implement me")
}

func (c Choice) Unfocus() tea.Cmd {
	return nil
}

func NewChoice(param, title string, options []string) *Choice {
	return &Choice{
		Param:   param,
		title:   title,
		options: options,
		choice:  -1,
	}
}
