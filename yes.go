package main

import (
	"fmt"
	"strings"
)

type Yes struct {
	output string
}

func NewYes(args []string) (Command, error) {
	output := ""

	if len(args) == 1 {
		output = "y"
	} else {
		output = strings.Join(args[1:], " ")
	}
	return Yes{output}, nil
}

func (y Yes) Run() int {
	for {
		fmt.Println(y.output)
	}
	return 0
}
