package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
)

func getWalkFunction(base string, data map[string]interface{}, target string) filepath.WalkFunc {
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
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		textContent := string(content)
		fileTemplate := template.New("filecontent")
		_, err = fileTemplate.Funcs(sprig.TxtFuncMap()).Parse(textContent)
		if err != nil {
			return err
		}

		newPath := filepath.Join(target, buf.String())
		os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
		newfile, err := os.OpenFile(newPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		fileTemplate.Execute(newfile, data)

		return nil
	}
}

func renderFolder(base string, data map[string]interface{}, target string) error {
	filepath.Walk(base, getWalkFunction(base, data, target))
	return nil
}
