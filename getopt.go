package main

import (
	"fmt"
	"strings"
)

const OPT_NO_ARG = 0
const OPT_REQUIRE_ARG = 1

const OPT_RET_END = -1
const OPT_RET_NO_ARG = -2

type Option struct {
	OptLong  string
	OptShort byte
	HasArg   int
}

func (o *Option) hasLong() bool {
	return o.OptLong != ""
}

func (o *Option) hasShort() bool {
	return o.OptShort != 0
}

func hasLongOptPrefix(arg string) bool {
	return strings.HasPrefix(arg, "--") && arg != "--"
}
func hasShortOptPrefix(arg string) bool {
	return strings.HasPrefix(arg, "-") && arg != "-"
}
func hasOptPrefix(arg string) bool {
	return strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-")
}

type Options struct {
	opts  []Option
	value string
	err   error
	opt   *Option

	//private
	argc int
	pos  int
}

func NewOptions(opt []Option) *Options {
	return &Options{opts: opt}
}

func (o *Options) Option() *Option {
	if o.err != nil {
		return nil
	} else {
		return o.opt
	}
}
func (o *Options) Value() string {
	if o.err != nil {
		return ""
	} else {
		return o.value
	}
}
func (o *Options) Err() error {
	return o.err
}

func (o *Options) GetOpts(args []string) bool {
	if o.err != nil {
		return false
	}

	for o.argc < len(args) {
		if o.pos == 0 {
			arg := args[o.argc]
			if hasLongOptPrefix(arg) {
				return o.getoptLong(args)
			} else if hasShortOptPrefix(arg) {
				return o.getoptShort(args)
			} else if !hasOptPrefix(arg) {
				o.value = arg
				o.opt = nil
				o.argc++
				return true
			} else {
				o.err = fmt.Errorf("Invalid option \"%s\"", arg)
				return false
			}
		} else {
			return o.getoptShort(args)
		}
	}
	return false
}

func (o *Options) getoptLong(args []string) bool {
	arg := args[o.argc][2:]
	tokens := strings.SplitN(arg, "=", 2)

	for _, opt := range o.opts {
		if opt.hasLong() && opt.OptLong == tokens[0] {
			//check args has an argument if the opt needs
			if opt.HasArg == OPT_REQUIRE_ARG {
				if len(tokens) == 2 {
					// --opt=value
					o.opt = &opt
					o.value = tokens[1]
					o.argc++
					return true
				} else if o.argc < len(args)-1 {
					// --opt value
					o.opt = &opt
					o.value = args[o.argc+1]
					o.argc += 2
					return true
				} else {
					o.err = fmt.Errorf("Option \"--%s\" needs an argument", arg)
					return false
				}
			} else {
				if len(tokens) == 2 {
					o.err = fmt.Errorf("Invalid argument %s for opt \"--%s\"", tokens[1], tokens[0])
					return false
				}
				o.opt = &opt
				o.value = ""
				o.argc++
				return true
			}
		}
	}

	//not found
	o.err = fmt.Errorf("Unknown option \"--%s\"", arg)
	return false
}
func (o *Options) getoptShort(args []string) bool {
	arg := args[o.argc]
	tokens := strings.SplitN(arg, "=", 2)

	//ignore option prefix
	for tokens[0][o.pos] == '-' && o.pos < len(tokens[0]) {
		o.pos++
	}

	if o.pos == len(tokens[0]) {
		o.err = fmt.Errorf("Unknown options %s", arg)
		return false
	}

	for _, opt := range o.opts {
		if opt.hasShort() && opt.OptShort == tokens[0][o.pos] {
			//found it
			//check if the opt is the last one if it needs an argument
			if opt.HasArg == OPT_REQUIRE_ARG {
				if o.pos != len(tokens[0])-1 {
					o.err = fmt.Errorf("Option \"-%c\" needs an argument", tokens[0][o.pos])
					return false
				}
				if len(tokens) == 2 {
					o.opt = &opt
					o.value = tokens[1]
					o.argc += 1
					o.pos = 0
					return true
				} else if o.argc != len(args)-1 {
					o.opt = &opt
					o.value = args[o.argc+1]
					o.argc += 2
					o.pos = 0
					return true
				} else {
					o.err = fmt.Errorf("Option \"-%c\" needs an argument", arg[o.pos])
					return false
				}
			} else {
				if o.pos == len(tokens[0])-1 && len(tokens) == 2 {
					o.err = fmt.Errorf("Invalid argument %s for opt \"-%c\"", tokens[1], tokens[0][o.pos])
					return false
				}
				o.opt = &opt
				o.value = ""
				o.pos++
				if o.pos == len(tokens[0]) {
					o.pos = 0
					o.argc++
				}
				return true
			}
		}
	}

	//not found
	o.err = fmt.Errorf("Unknown option \"-%c\"", arg[o.pos])
	return false
}
