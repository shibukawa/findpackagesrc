package findpackagesrc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFinder_FindSourcePath(t *testing.T) {
	f := Finder{
		goModPath: "/go/src/github.com/myproject",
		goPath:    "/go",
		sumPakcage: [][2]string{
			{"example.com/sub", "example.com/sub@1.0.0"},
		},
		replaces:   []replaceEntry{
			{
				Source: "example.com/a/b/",
				Target: "../b",
				Type:   LocalPath,
			},
			{
				Source: "example.com/a/",
				Target: "../../a",
				Type:   0,
			},
		},
	}
	testcases := []struct{
		name string
		source string
		expected string
		hasError bool
	}{
		{
			name:     "find by go.mod replace",
			source:   "example.com/a/b/c",
			expected:  "/go/src/github.com/b/c",
			hasError: false,
		},
		{
			name:     "find by go.sum",
			source:   "example.com/sub/subsub",
			expected: "/go/pkg/mod/example.com/sub@1.0.0/subsub",
			hasError: false,
		},
		{
			name:     "can't find",
			source:   "github.com/404/notfound",
			hasError: true,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.FindSourcePath(tt.source)
			if tt.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				if err != nil {
					return
				}
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}