package main

import (
	"fmt"
	"os"
	"path"
)

type cmdMainFunc func(args []string) int

var commands map[string]cmdMainFunc = map[string]cmdMainFunc{
	"yes":   YesMain,
	"head":  HeadMain,
	"sleep": SleepMain,
}

func busyboxUsage() {
	fmt.Println("busybox-go v0.0.1")
	fmt.Println()

	fmt.Println("Usage:")
	fmt.Println("\tbusybox-go [function [arguments]...]")
	fmt.Println("\tfunction [arguments]...")
	fmt.Println()

	fmt.Println("Currently defined functions:")

	for k := range commands {
		fmt.Printf("%s, ", k)
	}
	fmt.Println()
}

func run(args []string) int {
	name := path.Base(args[0])

	if name == "busybox" || name == "busybox-go" {
		if len(args) == 1 {
			busyboxUsage()
			return 0
		}
		return run(args[1:])
	}

	if cmd, ok := commands[name]; ok {
		return cmd(args)
	}

	fmt.Printf("%s: command not found\n", name)
	return -1
}

func main() {
	os.Exit(run(os.Args))
}
