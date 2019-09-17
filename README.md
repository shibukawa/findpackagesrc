# findpackagesrc

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