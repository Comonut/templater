package util

import (
	"strings"
	"testing"
)

func TestMerge(t *testing.T) {
	for _, test := range []struct {
		testid   string
		content1 string
		content2 string
		template string
		expected string
	}{
		{
			testid:   "simple merge",
			content1: "a\nb\nc",
			content2: "a\nY\nc",
			template: "a\nc",
			expected: "a\nb\nY\nc",
		},
		{
			testid:   "YAML with array of objects",
			content1: "kube:\n  params:\n    - param1:\n      abc: 1",
			content2: "kube:\n  params:\n    - param2:\n      abc: 2",
			template: "kube:\n  params:",
			expected: "kube:\n  params:\n    - param1:\n      abc: 1\n    - param2:\n      abc: 2",
		},
		{
			testid:   "YAML with two mergeable sections",
			content1: "apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\n\nnamePrefix: test-\n\nresources:\n  - first-cronjob.yaml\n\nimages:\n  - name: test-image\n\nconfigMapGenerator:\n  - {name: first-config, envs: [first-config.env]}\n\nconfigurations:\n  - kustomizeconfig/sealedsecretkind.yaml",
			content2: "apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\n\nnamePrefix: test-\n\nresources:\n  - second-cronjob.yaml\n\nimages:\n  - name: test-image\n\nconfigMapGenerator:\n  - {name: second-config, envs: [second-config.env]}\n\nconfigurations:\n  - kustomizeconfig/sealedsecretkind.yaml",
			template: "apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\n\nnamePrefix: test-\n\nresources:\n\nimages:\n  - name: test-image\n\nconfigMapGenerator:\n\nconfigurations:\n  - kustomizeconfig/sealedsecretkind.yaml",
			expected: "apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\n\nnamePrefix: test-\n\nresources:\n  - first-cronjob.yaml\n  - second-cronjob.yaml\n\nimages:\n  - name: test-image\n\nconfigMapGenerator:\n  - {name: first-config, envs: [first-config.env]}\n  - {name: second-config, envs: [second-config.env]}\n\nconfigurations:\n  - kustomizeconfig/sealedsecretkind.yaml",
		},
	} {
		firstReader := strings.NewReader(test.content1)
		secondReader := strings.NewReader(test.content2)
		templateReader := strings.NewReader(test.template)
		result := mergeContents(firstReader, secondReader, templateReader)
		if result != test.expected {
			t.Errorf("Failed `%s`\nexpected:\n%s\ngot:\n%s", test.testid, test.expected, result)
		}
	}

}
