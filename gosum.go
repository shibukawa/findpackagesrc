package findpackagesrc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

func parseGoSumFile(dir string) (map[string]string, error) {
	gosum := filepath.Join(dir, "go.sum")
	f, err := os.Open(gosum)
	if err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
				return nil, errors.New("use go modules")
			}
			return nil, errors.New("cannot read go.sum")
		}
		return nil, err
	}
	defer f.Close()
	return parseGoSum(f)
}

func parseGoSum(reader io.Reader) (map[string]string, error) {
	result := make(map[string]string)
	scr := bufio.NewScanner(reader)
	for scr.Scan() {
		stuff := strings.Fields(scr.Text())
		if len(stuff) != 3 {
			continue
		}
		if strings.HasSuffix(stuff[1], "/go.mod") {
			continue
		}
		encodedPath, err := encodeString(stuff[0])
		if err != nil {
			return nil, err
		}
		result[stuff[0]] = encodedPath + "@" + stuff[1]
	}
	if err := scr.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func convertGoSumMapToSlice(src map[string]string) [][2]string {
	result := make([][2]string, 0, len(src))
	for k, v := range src {
		result = append(result, [2]string{k, v})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i][0] > result[j][1]
	})
	return result
}

// copied from cmd/go/internal/module/module.go via https://github.com/Songmu/gocredits/blob/master/gocredits.go
func encodeString(s string) (encoding string, err error) {
	haveUpper := false
	for _, r := range s {
		if r == '!' || r >= utf8.RuneSelf {
			// This should be disallowed by CheckPath, but diagnose anyway.
			// The correctness of the encoding loop below depends on it.
			return "", fmt.Errorf("internal error: inconsistency in EncodePath")
		}
		if 'A' <= r && r <= 'Z' {
			haveUpper = true
		}
	}

	if !haveUpper {
		return s, nil
	}

	var buf []byte
	for _, r := range s {
		if 'A' <= r && r <= 'Z' {
			buf = append(buf, '!', byte(r+'a'-'A'))
		} else {
			buf = append(buf, byte(r))
		}
	}
	return string(buf), nil
}
