package main

import (
	"fmt"
	"os"
	"path"
)

type cmdFunc struct {
	newCmd func(args []string) (Command, error)
	usage  func()
}

var commands map[string]cmdFunc = initCmdList()

func initCmdList() map[string]cmdFunc {
	return map[string]cmdFunc{
		"yes":  cmdFunc{NewYes, UsageYes},
		"head": cmdFunc{NewHead, UsageHead},
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

	if c, ok := commands[name]; ok {
		cmd, err := c.newCmd(args)
		if err != nil {
			fmt.Println(err)
			c.usage()
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
