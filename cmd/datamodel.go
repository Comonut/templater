package cmd

import (
	"fmt"
	"github.com/Comonut/templater/components"
)

type datamodel map[string]interface{}

func (dm *datamodel) getParamType() (string, error) {
	t, ok := (*dm)["type"]
	if !ok {
		return "", fmt.Errorf("mising `type` value")
	}
	switch casted := t.(type) {
	case string:
		return casted, nil
	default:
		return "", fmt.Errorf("`type` value should be string")
	}
}

func (dm *datamodel) getParamName() (string, error) {
	name, ok := (*dm)["param"]
	if !ok {
		return "", fmt.Errorf("mising `param` value")
	}
	switch casted := name.(type) {
	case string:
		return casted, nil
	default:
		return "", fmt.Errorf("`param` value should be string")
	}
}

func (dm *datamodel) getParamTitle() (string, error) {
	name, ok := (*dm)["title"]
	if ok {
		switch casted := name.(type) {
		case string:
			return casted, nil
		default:
			return "", fmt.Errorf("`title` value should be string")
		}
	}
	return dm.getParamName()
}

func (dm *datamodel) getParamExample() (string, error) {
	name, ok := (*dm)["example"]
	if ok {
		switch casted := name.(type) {
		case string:
			return casted, nil
		default:
			return "", fmt.Errorf("`example` value should be string")
		}
	}
	return "", nil
}

func (dm *datamodel) getParamOptions() ([]string, error) {
	options, ok := (*dm)["options"]
	if ok {
		switch casted := options.(type) {
		case []interface{}:
			options := make([]string, len(casted))
			for i, _ := range casted {
				element, ok := casted[i].(string)
				if !ok {
					return []string{}, fmt.Errorf("`options` value should be a string array")
				}
				options[i] = element
			}
			return options, nil
		default:
			return []string{}, fmt.Errorf("`options` value should be a string array")
		}
	}
	return []string{}, nil
}

func (dm *datamodel) generateComponent() (components.Component, error) {
	t, err := dm.getParamType()
	if err != nil {
		return nil, err
	}
	switch t {
	case "textfield":
		return dm.generateTextField()
	case "choice":
		return dm.generateChoice()
	default:
		return nil, fmt.Errorf("unrecognized type %s", t)
	}
}

func (dm *datamodel) generateTextField() (*components.TextField, error) {
	param, err := dm.getParamName()
	if err != nil {
		return nil, err
	}
	title, err := dm.getParamTitle()
	if err != nil {
		return nil, err
	}

	example, err := dm.getParamExample()
	if err != nil {
		return nil, err
	}

	return components.NewTextField(param, title, example), nil
}

func (dm *datamodel) generateChoice() (*components.Choice, error) {
	param, err := dm.getParamName()
	if err != nil {
		return nil, err
	}
	title, err := dm.getParamTitle()
	if err != nil {
		return nil, err
	}

	options, err := dm.getParamOptions()
	if err != nil {
		return nil, err
	}

	return components.NewChoice(param, title, options), nil
}
