package findpackagesrc

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestFinder_FindSourcePath_GoMod(t *testing.T) {
	f := Finder{
		projectPath: "/home/go/src/github.com/myproject",
		goPath:      "/home/go",
		goRoot:      "/go",
		sumPakcage: [][2]string{
			{"example.com/sub", "example.com/sub@1.0.0"},
		},
		replaces: []replaceEntry{
			{
				Source: "example.com/a/b/",
				Target: "../b",
				Type:   LocalPath,
			},
			{
				Source: "example.com/a/",
				Target: "../../a",
				Type:   LocalPath,
			},
			{
				Source: "example.com/c/",
				Target: "example.com/sub",
				Type:   Module,
			},
		},
	}
	testcases := []struct {
		name     string
		source   string
		expected string
		hasError bool
	}{
		{
			name:     "find by go.mod replace",
			source:   "example.com/a/b/c",
			expected: "/home/go/src/github.com/b/c",
			hasError: false,
		},
		{
			name:     "find by go.sum",
			source:   "example.com/sub/subsub",
			expected: "/home/go/pkg/mod/example.com/sub@1.0.0/subsub",
			hasError: false,
		},
		{
			name:     "find by go.mod and go.sum",
			source:   "example.com/c/sub",
			expected: "/home/go/pkg/mod/example.com/sub@1.0.0/sub",
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

func must(path string, err error) string {
	if err != nil {
		panic(err)
	}
	return path
}

func TestFinder_FindSourcePath_traditional(t *testing.T) {
	gopath := must(filepath.Abs("./testdata/gopath_test"))
	projectPath := must(filepath.Abs("./testdata/vendor_test"))
	goroot := must(run("go", "env", "GOROOT"))
	f := Finder{
		projectPath: projectPath,
		goPath:      gopath,
		goRoot:      goroot,
		sumPakcage:  nil,
		replaces:    nil,
	}

	testcases := []struct {
		name        string
		goPath      string
		projectPath string
		source      string
		expected    string
		hasError    bool
	}{
		{
			name:     "GOPATH",
			source:   "example.com/my-package-in-gopath/sub",
			expected: filepath.Join(gopath, "src", "example.com", "my-package-in-gopath", "sub"),
			hasError: false,
		},
		{
			name:     "vendor",
			source:   "example.com/my-package-in-vendor/sub",
			expected: filepath.Join(projectPath, "vendor", "example.com", "my-package-in-vendor", "sub"),
			hasError: false,
		},
		{
			name:     "GOROOT",
			source:   "encoding/json",
			expected: filepath.Join(goroot, "src", "encoding", "json"),
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
