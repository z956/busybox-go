package main

import (
	"fmt"
	"os"
	"strconv"
)

type head struct {
	file    string
	isStdin bool
	num     int

	isLineMode bool
	isVerbose  bool
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

func (h *head) setArg(val string) error {
	//file
	fmt.Printf("set arg, val: %s\n", val)
	if _, err := os.Stat(val); os.IsNotExist(err) {
		return err
	} else {
		h.file = val
		h.isStdin = false
	}
	return nil
}
func (h *head) setOptArg(opt *Option, val string) error {
	var err error
	switch opt.OptShort {
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
		err = fmt.Errorf("Invalid option \"-%c\"", opt.OptShort)
		return err
	}
	return nil
}

func (h *head) run() int {
	fmt.Println("head command is running")
	return 0
}

func (h *head) usage() {
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

func HeadMain(args []string) int {
	h := head{"", true, 10, true, false}
	var err error

	opts := Options{
		opts: []Option{
			{"", 'n', OPT_REQUIRE_ARG},
			{"", 'c', OPT_REQUIRE_ARG},
			{"", 'v', OPT_NO_ARG},
			{"", 'q', OPT_NO_ARG},
		},
	}

	for err == nil {
		if i := opts.GetOpts(args[1:]); i == OPT_RET_END {
			break
		} else if i == OPT_RET_NO_ARG {
			err = h.setArg(opts.Value)
		} else {
			err = h.setOptArg(&opts.opts[i], opts.Value)
		}
	}

	if opts.Err == nil && err == nil {
		return h.run()
	}
	if opts.Err != nil {
		fmt.Println(opts.Err)
	} else if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	h.usage()
	return -1
}
