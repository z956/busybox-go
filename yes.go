package main

import (
	"fmt"
	"strings"
)

type Yes struct{}

func (m Yes) Run(args []string) {
	output := ""

	if len(args) == 1 {
		output = "y"
	} else {
		output = strings.Join(args[1:], " ")
	}
	for {
		fmt.Println(output)
	}
}
