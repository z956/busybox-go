package main

import (
	"fmt"
	"os"
	"path"
)

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
	var cmd Command

	switch path.Base(args[0]) {
	case "yes":
		cmd = Yes{}
	case "busybox", "busybox-go":
		if len(args) == 1 {
			busyboxUsage()
			return 0
		} else if args[1] == "busybox" || args[1] == "busybox-go" {
			busyboxUsage()
			return 0
		} else {
			return runCommand(args[1:])
		}
	default:
		fmt.Printf("%s: applet not found\n", args[0])
		return -1
	}

	return cmd.Run(args)
}

func main() {
	os.Exit(runCommand(os.Args))
}
