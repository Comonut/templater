package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type Choice struct {
	Param   string
	Title   string
	Options []string
	Choice  int
	Focused bool
}

func (c Choice) Update(msg tea.Msg) (Component, tea.Cmd) {
	if !c.Focused {
		return c, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if c.Choice == 0 {
				c.Choice = len(c.Options) - 1
				return c, nil
			}
			c.Choice -= 1
			return c, nil
		case "right":
			if c.Choice == len(c.Options)-1 {
				c.Choice = 0
				return c, nil
			}
			c.Choice += 1
			return c, nil
		}
	}
	return c, nil
}

func (c Choice) View() string {
	var title string
	if c.Focused {
		title = focusedStyle.Copy().Render(c.Title)
	} else {
		title = c.Title
	}

	var builder strings.Builder
	builder.WriteString(title)
	builder.WriteString("\n")
	for i, option := range c.Options {
		if i != 0 {
			builder.WriteString("\t")
		}
		var renderedOption string
		if i != c.Choice {
			renderedOption = blurredStyle.Copy().Render(option)
		} else {
			if c.Focused {
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
	if c.Choice == -1 {
		return ""
	}
	return c.Options[c.Choice]
}

func (c Choice) Focus() (Component, tea.Cmd) {
	c.Focused = true
	return c, nil
}

func (c Choice) Unfocus() (Component, tea.Cmd) {
	c.Focused = false
	return c, nil
}

func NewChoice(param, title string, options []string) *Choice {
	return &Choice{
		Param:   param,
		Title:   title,
		Options: options,
		Choice:  0,
	}
}
