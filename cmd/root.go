package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "templater <template folder> <output path>",
	Short: "Templater gives you a simpler way to create and run code repo templates",
	Args:  cobra.ExactValidArgs(2),
	Run:   run,
}

const ParamsFile = "params.yaml"
const TemplatesFolder = "templates"

var ValuesFile string

func init() {
	rootCmd.Flags().StringVarP(&ValuesFile, "values", "v", "", "File containing parameter values")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

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

func getValuesFromUser(datamodel []map[string]interface{}) map[string]interface{} {
	model, err := initModel(datamodel)
	if err != nil {
		fmt.Printf("Error initializing ui: %s\n", err)
		os.Exit(1)
	}

	if err := tea.NewProgram(model).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	data := model.retrieveData()
	return data
}

func run(cmd *cobra.Command, args []string) {

	template := args[0]
	target := args[1]

	params, err := getDataModel(filepath.Join(template, ParamsFile))
	if err != nil {
		fmt.Printf("Error loading data model: %s\n", err)
		os.Exit(1)
	}
	var data map[string]interface{}
	if ValuesFile == "" {
		data = getValuesFromUser(params)
	} else {
		//TODO parse values file
		data = map[string]interface{}{}
	}

	err = renderFolder(filepath.Join(template, TemplatesFolder), data, target)
	if err != nil {
		fmt.Printf("Error rendering templates: %s\n", err)
		os.Exit(1)
	}

}
