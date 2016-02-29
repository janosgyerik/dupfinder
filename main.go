package main

import (
	"fmt"
	"github.com/janosgyerik/go-practice/dupfinder/dupfinder"
)

func main() {
	fmt.Println(dupfinder.Compare("/tmp/main.go", "/tmp/main.go"))
	fmt.Println(dupfinder.Compare("/tmp/main.go", "/tmp/a"))
	fmt.Println(dupfinder.Compare("/tmp/a", "/tmp/main.go"))
}

