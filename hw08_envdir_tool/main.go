package main

import (
	"fmt"
	"os"
)

func main() {
	// fmt.Println(os.Args)

	environment, err := ReadDir(os.Args[1])
	// fmt.Println(environment)
	if err != nil {
		fmt.Println(err)
	}

	RunCmd(os.Args[2:], environment)
}
