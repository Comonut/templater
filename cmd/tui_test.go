package cmd

import (
	"github.com/Comonut/templater/components"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

func equalModels(this, other model) (bool, string) {
	if this.focusIndex != other.focusIndex {
		return false, "cursor"
	}

	if this.submitted != other.submitted {
		return false, "submitted"
	}

	if len(this.inputs) != len(other.inputs) {
		return false, "inputs size"
	}

	for i := range this.inputs {
		switch thisCasted := this.inputs[i].(type) {
		case components.TextField:
			otherCasted, ok := other.inputs[i].(*components.TextField)
			if !ok {
				return false, "component types don't match"
			}
			if thisCasted.Param != otherCasted.Param {
				return false, "textfield param doesn't match"
			}
			if thisCasted.Input.Placeholder != otherCasted.Input.Placeholder {
				return false, "textfield placeholder doesn't match"
			}
		case components.Choice:
			otherCasted, ok := other.inputs[i].(*components.Choice)
			if !ok {
				return false, "component types don't match"
			}
			if thisCasted.Param != otherCasted.Param {
				return false, "choice param doesn't match"
			}
			if thisCasted.Choice != otherCasted.Choice {
				return false, "choice selection doesn't match"
			}

		}
	}

	return true, ""
}

func textField(placeholder string, focused bool) textinput.Model {
	t := textinput.NewModel()
	t.Placeholder = placeholder
	if focused {
		t.Focus()
	}
	return t
}

func TestGenerateModel(t *testing.T) {
	for _, test := range []struct {
		name      string
		datamodel []datamodel
		expected  model
	}{
		{
			name: "simple two text fields",
			datamodel: []datamodel{
				{
					"param": "name",
					"type":  "textfield",
				},
				{
					"param":   "age",
					"type":    "textfield",
					"title":   "Age",
					"example": "23",
				},
			},
			expected: model{
				focusIndex: 0,
				inputs: []components.Component{
					components.NewTextField("name", "name", ""),
					components.NewTextField("age", "Age", "23"),
				},
				submitted: false,
			},
		},
	} {
		generated, err := initModel(test.datamodel)
		if err != nil {
			t.Errorf("Failed %s - %s", test.name, err)
		}

		if ok, diff := equalModels(generated, test.expected); !ok {
			t.Errorf("Failed %s - %s doesn't match expected", test.name, diff)
		}
	}

}

func TestUpdateModel(t *testing.T) {
	for _, test := range []struct {
		name     string
		current  model
		message  tea.Msg
		expected model
	}{
		{
			name: "down button updates index",
			current: model{
				focusIndex: 0,
				inputs: []components.Component{
					components.NewTextField("name", "name", ""),
					components.NewTextField("age", "Age", "23"),
				},
			},
			message: tea.KeyMsg{
				Type: tea.KeyDown,
			},
			expected: model{
				focusIndex: 1,
				inputs: []components.Component{
					components.NewTextField("name", "name", ""),
					components.NewTextField("age", "Age", "23"),
				},
			},
		},
		{
			name: "right on choice updates choice",
			current: model{
				focusIndex: 0,
				inputs: []components.Component{
					&components.Choice{
						Param:   "test",
						Title:   "test",
						Options: []string{"one", "two"},
						Choice:  0,
						Focused: true,
					},
				},
			},
			message: tea.KeyMsg{
				Type: tea.KeyRight,
			},
			expected: model{
				focusIndex: 0,
				inputs: []components.Component{
					&components.Choice{
						Param:   "test",
						Title:   "test",
						Options: []string{"one", "two"},
						Choice:  1,
						Focused: true,
					},
				},
			},
		},
		{
			name: "right on choice doesn't update choice when unfocused",
			current: model{
				focusIndex: 1,
				inputs: []components.Component{
					&components.Choice{
						Param:   "test",
						Title:   "test",
						Options: []string{"one", "two"},
						Choice:  0,
						Focused: false,
					},
				},
			},
			message: tea.KeyMsg{
				Type: tea.KeyRight,
			},
			expected: model{
				focusIndex: 1,
				inputs: []components.Component{
					&components.Choice{
						Param:   "test",
						Title:   "test",
						Options: []string{"one", "two"},
						Choice:  0,
						Focused: false,
					},
				},
			},
		},
	} {
		updated, _ := test.current.Update(test.message)
		switch casted := updated.(type) {
		case model:
			if ok, diff := equalModels(casted, test.expected); !ok {
				t.Errorf("Failed %s - %s doesn't match expected", test.name, diff)
			}
		default:
			t.Errorf("Failed %s - Update didn't return a `model` struct", test.name)
		}
	}
}
