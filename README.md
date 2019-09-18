# findpackagesrc

[![GoDoc](https://godoc.org/github.com/shibukawa/findpackagesrc?status.svg)](https://godoc.org/github.com/shibukawa/findpackagesrc)

Golang package that find original source path

```go
package main

import (
    "fmt"
    "github.com/shibukawa/findpackagesrc"
)

func main() {
    finder, err := findpackagesrc.NewFinder(findpackagesrc.Option{})
    if err != nil {
        panic(err)
    }
    path, err := finder.FindSourcePath("github.com/stretchr/testify")
    if err != nil {
        panic(err)
    }
    fmt.Println(path)
    // /home/yourname/go/pkg/mod/github.com/stretchr/testify@v1.4.0
}
```

It searches in the following order (I don't know the order is as same as Go's rule):

1. Replace defined in go.mod
2. Downloaded by go mod in go.sum
3. vendor folder
4. $GOPATH
5. $GOROOT (for go standard library)

## License

Apache 2