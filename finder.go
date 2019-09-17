// findpackagesrc package provides API to detect original source of go package
package findpackagesrc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Finder has method returns package source
type Finder struct {
	goModPath  string
	goPath     string
	sumPakcage [][2]string
	replaces   []replaceEntry
}

// Option specify Finder's option
//
// Default value of PkgPath is a nearest parent folder that contains go.mod.
//
// Default value of GoPath is a $GOPATH.
type Option struct {
	PkgPath string
	GoPath  string
}

// NewFinder is a constructor function of Finder
func NewFinder(option Option) (*Finder, error) {
	if option.PkgPath == "" {
		path, err := filepath.Abs(".")
		if err != nil {
			return nil, fmt.Errorf("Unknown error when detect current folder: %w", err)
		}
		for {
			_, err = os.Stat(filepath.Join(path, "go.mod"))
			if err == nil {
				option.PkgPath = path
				break
			}
			if os.IsNotExist(err) {
				parent := filepath.Dir(path)
				if parent == path {
					return nil, fmt.Errorf("Can't detect go package (go.mod) directory or use option.PkgPath: %w", err)
				}
				path = parent
			} else {
				return nil, fmt.Errorf("Unknown error when detect current folder: %w", err)
			}
		}
	}
	if option.GoPath == "" {
		goroot, err := run("go", "env", "GOPATH")
		if err != nil {
			return nil, fmt.Errorf("Can't detect $GOPATH. go is not in PATH or use option.GoPath: %w", err)
		}
		option.GoPath = goroot
	}
	fmt.Println(option)
	packages, err := parseGoSumFile(option.PkgPath)
	fmt.Println(packages)
	if err != nil {
		return nil, fmt.Errorf("Error detected during parsing go.sum: %w", err)
	}
	replaces, err := parseGoModFile(option.PkgPath)
	fmt.Println(replaces)
	if err != nil {
		return nil, fmt.Errorf("Error detected during parsing go.mod: %w", err)
	}
	return &Finder{
		goModPath:  option.PkgPath,
		goPath:     option.GoPath,
		sumPakcage: convertGoSumMapToSlice(packages),
		replaces:   convertMapToSlice(replaces),
	}, nil
}

// FindSourcePath returns source folder of specified package
func (f Finder) FindSourcePath(pkg string) (string, error) {
	if !strings.HasSuffix(pkg, "/") {
		pkg += "/"
	}
	for _, replace := range f.replaces {
		if strings.HasPrefix(pkg, replace.Source) {
			rest := pkg[len(replace.Source):]
			return filepath.Join(f.goModPath, replace.Target, rest), nil
		}
	}
	for _, mod := range f.sumPakcage {
		if strings.HasPrefix(pkg, mod[0]) {
			rest := pkg[len(mod[0]):]
			return filepath.Join(f.goPath, "pkg", "mod", mod[1], rest), nil
		}
	}
	return "", errors.New("can't find packages")
}

