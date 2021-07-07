package util

import (
	"io"
)

func MergeContents(first, second, template io.Reader) string {
	merged, _ := Merge(first, template, second, true, "old", "new")
	result, _ := io.ReadAll(merged.Result)

	return string(result)
}
