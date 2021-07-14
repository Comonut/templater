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
	RunE:  run,
}

const ParamsFile = "params.yaml"
const TemplatesFolder = "templates"

var ValuesFile string
var WriteModeString string

func init() {
	rootCmd.Flags().StringVarP(&ValuesFile, "values", "v", "", "File containing parameter values")
	rootCmd.Flags().StringVarP(&WriteModeString, "mode", "m", "replace", "Output writing mode - one of append, ignore, replace, merge")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getDataModel(path string) ([]datamodel, error) {

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var values []datamodel
	err = yaml.Unmarshal(yamlFile, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func getValuesFromUser(datamodel []datamodel) map[string]interface{} {
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

func getValuesFromFile() (map[string]interface{}, error) {
	yamlFile, err := ioutil.ReadFile(ValuesFile)
	if err != nil {
		return nil, err
	}

	var values map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func run(cmd *cobra.Command, args []string) error {

	template := args[0]
	target := args[1]

	params, err := getDataModel(filepath.Join(template, ParamsFile))
	if err != nil {
		return err
	}
	var data map[string]interface{}
	if ValuesFile == "" {
		data = getValuesFromUser(params)
	} else {
		data, err = getValuesFromFile()
		if err != nil {
			return err
		}
	}

	err = renderFolder(filepath.Join(template, TemplatesFolder), data, target, WriteModeString)
	return err
}
