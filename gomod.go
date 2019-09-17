package findpackagesrc

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ModType int

const (
	Module ModType = iota + 1
	LocalPath
)

type replaceEntry struct {
	Source string
	Target string
	Type   ModType
}

func parseGoModFile(dir string) (map[string]replaceEntry, error) {
	gomod := filepath.Join(dir, "go.mod")
	f, err := os.Open(gomod)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("cannot read go.mod")
		}
		return nil, err
	}
	defer f.Close()
	return parseGoMod(f)
}

func at(s []string, pos int) string {
	if len(s) > pos {
		return s[pos]
	} else {
		return ""
	}
}

func convertMapToSlice(src map[string]replaceEntry) []replaceEntry {
	result := make([]replaceEntry, 0, len(src))
	for _, entry := range src {
		result = append(result, entry)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Source > result[j].Source
	})
	return result
}

func parseGoMod(r io.Reader) (map[string]replaceEntry, error) {
	result := make(map[string]replaceEntry)
	scr := bufio.NewScanner(r)
	inReplace := false
	for scr.Scan() {
		stuff := strings.Fields(scr.Text())
		var source string
		var target string
		if inReplace {
			if at(stuff, 0) == ")" {
				inReplace = false
				continue
			} else if at(stuff, 2) == "=>" {
				source = at(stuff, 0)
				target = at(stuff, 3)
			} else if at(stuff, 1) == "=>" {
				source = at(stuff, 0)
				target = at(stuff, 2)
			}
		} else if at(stuff, 0) == "replace" {
			if at(stuff, 1) == "(" {
				inReplace = true
				continue
			} else if at(stuff, 3) == "=>" {
				source = at(stuff, 1)
				target = at(stuff, 4)
			} else if at(stuff, 2) == "=>" {
				source = at(stuff, 1)
				target = at(stuff, 3)
			}
		}
		if source != "" && target != "" {
			if !strings.HasSuffix(source, "/") {
				source += "/"
			}
			if strings.HasPrefix(target, ".") || strings.HasPrefix(target, "/") {
				result[source] = replaceEntry{
					Source: source,
					Target: target,
					Type:   LocalPath,
				}
			} else {
				result[source] = replaceEntry{
					Source: source,
					Target: target,
					Type:   Module,
				}
			}
		}
	}
	if err := scr.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
