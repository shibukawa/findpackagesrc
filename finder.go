// findpackagesrc package provides API to detect original source of go package
package findpackagesrc

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Finder has method returns package source
type Finder struct {
	projectPath string
	goPath      string
	goRoot      string
	sumPakcage  [][2]string
	replaces    []replaceEntry
}

// Option specify Finder's option
//
// Default value of ProjectPath is a is a nearest parent folder that contains go.mod.
// If go.mod is not found, it becomes current folder.
//
// Default value of GoPath is a $GOPATH.
//
// Default value of GoRoot is a $GOROOT.
type Option struct {
	ProjectPath string
	GoPath      string
	GoRoot      string
}

// NewFinder is a constructor function of Finder
func NewFinder(options ...Option) (*Finder, error) {
	var option Option
	if len(options) > 0 {
		option = options[0]
	}
	if option.ProjectPath == "" {
		currentPath, err := filepath.Abs(".")
		if err != nil {
			return nil, fmt.Errorf("Unknown error when detect current folder: %w", err)
		}
		path := currentPath
		for {
			_, err = os.Stat(filepath.Join(path, "go.mod"))
			if err == nil {
				option.ProjectPath = path
				break
			} else if os.IsNotExist(err) {
				parent := filepath.Dir(path)
				if parent == path {
					option.ProjectPath = currentPath
					break
				}
				path = parent
			} else {
				return nil, fmt.Errorf("Unknown error when detect project root: %w", err)
			}
		}
	}

	if option.GoPath == "" {
		gopath, err := run("go", "env", "GOPATH")
		if err != nil {
			return nil, fmt.Errorf("Can't detect $GOPATH. go is not in PATH or use option.GoPath: %w", err)
		}
		option.GoPath = gopath
	}

	if option.GoRoot == "" {
		goroot, err := run("go", "env", "GOROOT")
		if err != nil {
			return nil, fmt.Errorf("Can't detect $GOROOT. go is not in PATH or use option.GoRoot: %w", err)
		}
		option.GoRoot = goroot
	}

	packages, err := parseGoSumFile(option.ProjectPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("Error detected during parsing go.sum: %w", err)
	}

	replaces, err := parseGoModFile(option.ProjectPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("Error detected during parsing go.mod: %w", err)
	}

	return &Finder{
		projectPath: option.ProjectPath,
		goPath:      option.GoPath,
		goRoot:      option.GoRoot,
		sumPakcage:  convertGoSumMapToSlice(packages),
		replaces:    convertMapToSlice(replaces),
	}, nil
}

// FindSourcePath returns source folder of specified package
//
// This method searches package in the following order:
//
//   1. Replace defined in go.mod
//   2. Downloaded by go mod in go.sum
//   3. vendor folder
//   4. $GOPATH
//   5. $GOROOT (for go standard library)
func (f Finder) FindSourcePath(pkg string) (string, error) {
	if !strings.HasSuffix(pkg, "/") {
		pkg += "/"
	}
	for _, replace := range f.replaces {
		if strings.HasPrefix(pkg, replace.Source) {
			rest := pkg[len(replace.Source):]
			if replace.Type == LocalPath {
				return filepath.Join(f.projectPath, replace.Target, rest), nil
			}
			pkg = path.Join(replace.Target, rest)
			if !strings.HasSuffix(pkg, "/") {
				pkg += "/"
			}
			break
		}
	}
	for _, mod := range f.sumPakcage {
		if strings.HasPrefix(pkg, mod[0]) {
			rest := pkg[len(mod[0]):]
			return filepath.Join(f.goPath, "pkg", "mod", mod[1], rest), nil
		}
	}

	vendorPath := filepath.Join(f.projectPath, "vendor", pkg)
	if s, err := os.Stat(vendorPath); err == nil && s.IsDir() {
		return vendorPath, nil
	}

	gopathPath := filepath.Join(f.goPath, "src", pkg)
	if s, err := os.Stat(gopathPath); err == nil && s.IsDir() {
		return gopathPath, nil
	}

	gorootPath := filepath.Join(f.goRoot, "src", pkg)
	if s, err := os.Stat(gorootPath); err == nil && s.IsDir() {
		return gorootPath, nil
	}

	return "", errors.New("can't find packages")
}
