package findpackagesrc

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var (
	singleLineReplace         = `replace example.com/some/dependency => example.com/some/dependency v1.2.3`
	singleLineReplaceVersion  = `replace example.com/some/dependency v1.2.5 => example.com/some/dependency v1.2.3`
	singleLineReplaceFork     = `replace example.com/some/dependency => example.com/some/dependency-fork v1.2.3`
	singleLineReplaceRelative = `replace example.com/some/dependency => ../relative`
	singleLineReplaceAbsolute = `replace example.com/some/dependency => /root/module/path`
	multiLineReplaceAbsolute  = `replace (
    example.com/some/dependency => /root/module/path
)`
)

func Test_parseGoMod(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		expectedTarget string
		expectedType   ModType
	}{
		{
			name:           "single line",
			source:         singleLineReplace,
			expectedTarget: "example.com/some/dependency",
			expectedType:   Module,
		},
		{
			name:           "single line with version",
			source:         singleLineReplaceVersion,
			expectedTarget: "example.com/some/dependency",
			expectedType:   Module,
		},
		{
			name:           "single line with other repository",
			source:         singleLineReplaceFork,
			expectedTarget: "example.com/some/dependency-fork",
			expectedType:   Module,
		},
		{
			name:           "single line with local path",
			source:         singleLineReplaceRelative,
			expectedTarget: "../relative",
			expectedType:   LocalPath,
		},
		{
			name:           "single line",
			source:         singleLineReplaceAbsolute,
			expectedTarget: "/root/module/path",
			expectedType:   LocalPath,
		},
		{
			name:           "multi line",
			source:         multiLineReplaceAbsolute,
			expectedTarget: "/root/module/path",
			expectedType:   LocalPath,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modMap, err := parseGoMod(strings.NewReader(tt.source))
			assert.Nil(t, err)
			if err != nil {
				return
			}
			result, ok := modMap["example.com/some/dependency/"]
			assert.True(t, ok)
			if !ok {
				return
			}
			assert.Equal(t, tt.expectedType, result.Type)
			assert.Equal(t, tt.expectedTarget, result.Target)
		})
	}
}

func Test_convertMapToSlice(t *testing.T) {
	source := map[string]replaceEntry{
		"sample.com/a/": {
			Source: "sample.com/a/",
			Target: "../a",
			Type:   LocalPath,
		},
		"sample.com/a/b/": {
			Source: "sample.com/a/b/",
			Target: "../a/b",
			Type:   LocalPath,
		},
		"sample.com/a/c/": {
			Source: "sample.com/a/c/",
			Target: "../a/c",
			Type:   LocalPath,
		},
	}
	result := convertMapToSlice(source)
	assert.Equal(t, 3, len(result))
	if 3 != len(result) {
		return
	}
	assert.Equal(t, "sample.com/a/c/", result[0].Source)
	assert.Equal(t, "sample.com/a/b/", result[1].Source)
	assert.Equal(t, "sample.com/a/", result[2].Source)
}
