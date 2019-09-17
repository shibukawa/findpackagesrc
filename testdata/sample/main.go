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
}