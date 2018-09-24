package main

import "fmt"

type Head struct{}

func NewHead(args []string) (Command, error) {
	return nil, fmt.Errorf("head is not implement yet\n")
}

func UsageHead() {
	fmt.Println("Usage: head [OPTIONS] [FILE]...")
	fmt.Println()

	fmt.Println("Print first 10 lines of each FILE (or stdin) to stdout.")
	fmt.Println("With more than one FILE, precede each with a filename header.")
	fmt.Println()

	fmt.Printf("\t-n N[kbm]\tPrint first N lines\n")
	fmt.Printf("\t-n -N[kbm]\tPrint all except N last ilnes\n")
	fmt.Printf("\t-c [-]N[kbm]\tPrint first N bytes\n")
	fmt.Printf("\t-q\t\tNever print headers\n")
	fmt.Printf("\t-v\t\tAlways print headers\n")
}

func (h Head) Run() int {
	fmt.Println("head running")
	return 0
}
