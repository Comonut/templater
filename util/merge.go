package util

import (
	"io"
	"strings"
)

func MergeContents(first, second, template string) string {
	firstReader := strings.NewReader(first)
	secondReader := strings.NewReader(second)
	templateReader := strings.NewReader(template)
	merged, _ := Merge(firstReader, templateReader, secondReader, true, "old", "new")
	result, _ := io.ReadAll(merged.Result)

	return string(result)
}
