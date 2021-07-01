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
	focused bool
}

func (c Choice) Update(msg tea.Msg) (Component, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if c.choice == 0 {
				c.choice = len(c.options) - 1
				return c, nil
			}
			c.choice -= 1
			return c, nil
		case "right":
			if c.choice == len(c.options)-1 {
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
	var title string
	if c.focused {
		title = focusedStyle.Copy().Render(c.title)
	} else {
		title = c.title
	}

	var builder strings.Builder
	builder.WriteString(title)
	builder.WriteString("\n")
	for i, option := range c.options {
		if i != 0 {
			builder.WriteString("\t")
		}
		var renderedOption string
		if i != c.choice {
			renderedOption = blurredStyle.Copy().Render(option)
		} else {
			if c.focused {
				renderedOption = focusedStyle.Copy().Render(option)
			} else {
				renderedOption = option
			}
		}

		builder.WriteString(renderedOption)
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

func (c Choice) Focus() (Component, tea.Cmd) {
	c.focused = true
	return c, nil
}

func (c Choice) Unfocus() (Component, tea.Cmd) {
	c.focused = false
	return c, nil
}

func NewChoice(param, title string, options []string) *Choice {
	return &Choice{
		Param:   param,
		title:   title,
		options: options,
		choice:  0,
	}
}
