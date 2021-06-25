package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v2"
)

func getDataModel(path string) ([]map[string]interface{}, error) {

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var values []map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func main() {
	args := os.Args
	if len(args) != 3 {
		println("Wrong number of arguements")
		println("Use templater <template path> <output path>")
		os.Exit(1)
	}
	template := args[1]
	target := args[2]

	values, err := getDataModel(filepath.Join(template, "values.yaml"))
	if err != nil {
		fmt.Printf("Error loading data model: %s\n", err)
		os.Exit(1)
	}

	model, err := initModel(values)
	if err != nil {
		fmt.Printf("Error initializing ui: %s\n", err)
		os.Exit(1)
	}

	if err := tea.NewProgram(model).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	data := model.retrieveData()
	renderFolder(filepath.Join(template, "templates"), data, target)

}
