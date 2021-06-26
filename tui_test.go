package main

import (
	"fmt"
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
		if this.inputs[i].Placeholder != other.inputs[i].Placeholder {
			return false, fmt.Sprintf("placeholder for input %d", i)
		}
		if this.inputs[i].Focused() != other.inputs[i].Focused() {
			return false, fmt.Sprintf("focus for input %d", i)
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
		datamodel []map[string]interface{}
		expected  model
	}{
		{
			name: "simple two text fields",
			datamodel: []map[string]interface{}{
				{
					"param": "name",
				},
				{
					"param": "age",
				},
			},
			expected: model{
				focusIndex: 0,
				inputs: []textinput.Model{
					textField("name", true),
					textField("age", false),
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
				inputs: []textinput.Model{
					textField("name", true),
					textField("age", false),
				},
			},
			message: tea.KeyMsg{
				Type: tea.KeyDown,
			},
			expected: model{
				focusIndex: 1,
				inputs: []textinput.Model{
					textField("name", false),
					textField("age", true),
				},
			},
		},
		{
			name: "down on submit goes to first field",
			current: model{
				focusIndex: 2,
				inputs: []textinput.Model{
					textField("name", false),
					textField("age", false),
				},
			},
			message: tea.KeyMsg{
				Type: tea.KeyDown,
			},
			expected: model{
				focusIndex: 0,
				inputs: []textinput.Model{
					textField("name", true),
					textField("age", false),
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
