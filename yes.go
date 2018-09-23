package main

import (
	"fmt"
	"strings"
)

type Yes struct{}

func NewYes() Command {
	return Yes{}
}

func (m Yes) Run(args []string) int {
	output := ""

	if len(args) == 1 {
		output = "y"
	} else {
		output = strings.Join(args[1:], " ")
	}
	for {
		fmt.Println(output)
	}
	return 0
}
