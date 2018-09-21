package main

import (
	"os"
	"path"
)

func getCommand(name string) Command {
	basename := path.Base(name)
	switch basename {
	case "yes":
		return Yes{}
	default:
		return nil
	}
}

func main() {
	s := getCommand(os.Args[0])
	if s != nil {
		s.Run(os.Args)
	}
}
