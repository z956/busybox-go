package main

type Command interface {
	Run(args []string)
}
