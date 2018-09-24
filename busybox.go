package main

import (
	"fmt"
	"os"
	"path"
)

type newCmdFunc func(args []string) (Command, error)

var commands map[string]newCmdFunc = initCmdList()

func initCmdList() map[string]newCmdFunc {
	return map[string]newCmdFunc{
		"yes": NewYes,
	}
}

func busyboxUsage() {
	fmt.Println("busybox-go v0.0.1")
	fmt.Println()

	fmt.Println("Usage:")
	fmt.Println("\tbusybox-go [function [arguments]...]")
	fmt.Println("\tfunction [arguments]...")
	fmt.Println()

	fmt.Println("Currently defined functions:")
	fmt.Println("\tyes")
}

func runCommand(args []string) int {
	name := path.Base(args[0])

	if name == "busybox" || name == "busybox-go" {
		if len(args) == 1 {
			busyboxUsage()
			return 0
		}
		return runCommand(args[1:])
	}

	if f, ok := commands[name]; ok {
		cmd, err := f(args)
		if err != nil {
			fmt.Println(err)
			return -1
		}
		return cmd.Run()
	}

	fmt.Printf("%s: command not found\n", name)
	return -1
}

func main() {
	os.Exit(runCommand(os.Args))
}
