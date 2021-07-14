package cmd

import (
	"bytes"
	"fmt"
	"github.com/Comonut/templater/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
)

const (
	Append  = "append"
	Replace = "replace"
	Ignore  = "ignore"
	Merge   = "merge"
)

func appendFile(templateContent, targetPath string, data map[string]interface{}) error {
	file, err := os.OpenFile(targetPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err
	}

	fileTemplate := template.New("filecontent")
	_, err = fileTemplate.Funcs(sprig.TxtFuncMap()).Parse(templateContent)

	fileTemplate.Execute(file, data)
	return nil
}

func replaceFile(templateContent, targetPath string, data map[string]interface{}) error {
	file, err := os.OpenFile(targetPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer file.Close()
	if err != nil {
		return err
	}

	fileTemplate := template.New("filecontent")
	_, err = fileTemplate.Funcs(sprig.TxtFuncMap()).Parse(templateContent)

	file.Truncate(0)
	fileTemplate.Execute(file, data)
	return nil
}

func ignoreFile(templateContent, targetPath string, data map[string]interface{}) error {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) { //write only if target doesn't exist
		return replaceFile(templateContent, targetPath, data)
	}
	return nil
}

func mergeFile(templateContent, targetPath string, data map[string]interface{}) error {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return replaceFile(templateContent, targetPath, data)
	}
	file, err := os.OpenFile(targetPath, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fileTemplate := template.New("filecontent")
	_, err = fileTemplate.Funcs(sprig.TxtFuncMap()).Parse(templateContent)

	renderedContent := &bytes.Buffer{}
	fileTemplate.Execute(renderedContent, data)

	currentContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return err
	}

	mergedContent := util.MergeContents(string(currentContent), renderedContent.String(), templateContent)

	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.WriteString(mergedContent)
	if err != nil {
		return err
	}

	return nil
}

func readTemplate(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	textContent := string(content)
	return textContent, nil
}

func getWalkFunction(base string, data map[string]interface{}, target string, writeFunc func(template, targetPath string, values map[string]interface{}) error) filepath.WalkFunc {
	baseLen := len(base)
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		buf := &bytes.Buffer{}
		relativePath := path[baseLen:]
		t := template.New("filename")
		_, err = t.Funcs(sprig.TxtFuncMap()).Parse(relativePath)
		if err != nil {
			return err
		}
		err = t.Execute(buf, data)
		if err != nil {
			return err
		}

		templateContent, err := readTemplate(path)
		if err != nil {
			return err
		}

		newPath := filepath.Join(target, buf.String())
		os.MkdirAll(filepath.Dir(newPath), os.ModePerm)

		return writeFunc(templateContent, newPath, data)
	}
}

func renderFolder(base string, data map[string]interface{}, target string, writingMode string) error {
	var writeFunc func(template, targetPath string, values map[string]interface{}) error
	switch writingMode {
	case Append:
		writeFunc = appendFile
	case Replace:
		writeFunc = replaceFile
	case Ignore:
		writeFunc = ignoreFile
	case Merge:
		writeFunc = mergeFile
	default:
		return fmt.Errorf("unknown writeMode") //should be unreachable
	}
	filepath.Walk(base, getWalkFunction(base, data, target, writeFunc))
	return nil
}
