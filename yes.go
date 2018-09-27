package main

import (
	"fmt"
	"strings"
)

func YesMain(args []string) int {
	output := ""
	args = args[1:]

	if len(args) == 0 {
		output = "y"
	} else {
		output = strings.Join(args, " ")
	}

	for {
		fmt.Println(output)
	}
	return 0
}
