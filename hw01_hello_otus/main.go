package main

import (
	"fmt"
	"golang.org/x/example/stringutil"
)

func main() {
	var text string = "Hello, OTUS!"
	fmt.Println(stringutil.Reverse(text))
}
