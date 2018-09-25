package main

import (
	"fmt"
	"os"
	"strconv"
)

type Head struct {
	file    string
	isStdin bool
	num     int

	isLineMode bool
	isVerbose  bool
}

var headOpts = initOpt()

func initOpt() []Option {
	return []Option{
		{"", 'n', OPT_REQUIRE_ARG},
		{"", 'c', OPT_REQUIRE_ARG},
		{"", 'v', OPT_NO_ARG},
		{"", 'q', OPT_NO_ARG},
	}
}

func parseSizeOpt(val string) (int, error) {
	//val: 10
	//val: 10k, 10K (1024)
	//val: 10b, 10B (512)
	//val: 10m, 10M (1024 * 1024)
	//others are invalid
	last := val[len(val)-1]
	mul := 1
	switch last {
	case 'k', 'K':
		val = val[:len(val)-1]
		mul = 1024
	case 'b', 'B':
		val = val[:len(val)-1]
		mul = 512
	case 'm', 'M':
		val = val[:len(val)-1]
		mul = 1024 * 1024
	default:
	}

	if size, err := strconv.Atoi(val); err != nil {
		return 0, fmt.Errorf("Invalid size value: %s", val)
	} else {
		size *= mul
		return size, nil
	}
}
func parseResult(optlong string, optshort byte, val string, userdata interface{}) error {
	h := userdata.(*Head)
	var err error
	switch optshort {
	case 'n':
		if h.num, err = parseSizeOpt(val); err != nil {
			return err
		}
		h.isLineMode = true
	case 'c':
		if h.num, err = parseSizeOpt(val); err != nil {
			return err
		}
		h.isLineMode = false
	case 'v':
		h.isVerbose = true
	case 'q':
		h.isVerbose = false
	default:
		//file
		if _, err := os.Stat(val); os.IsNotExist(err) {
			return err
		} else {
			h.file = val
			h.isStdin = false
		}
	}
	return nil
}

func NewHead(args []string) (Command, error) {
	h := Head{"", true, 10, true, false}

	if err := OptParse(args[1:], headOpts, parseResult, &h); err != nil {
		return nil, err
	}

	return h, nil
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
	fmt.Printf("head running, head: %v\n", h)
	return 0
}
