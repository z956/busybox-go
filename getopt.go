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
	Value string
	Err   error

	//private
	argc int
	pos  int
}

func (o *Options) GetOpts(args []string) int {
	if o.Err != nil {
		return OPT_RET_END
	}

	argslen := len(args)

	for o.argc < argslen {
		if o.pos == len(args[o.argc]) {
			o.pos = 0
			o.argc++
			continue
		}

		if o.pos == 0 {
			arg := args[o.argc]
			if hasLongOptPrefix(args[o.argc]) {
				return o.getoptLong(args)
			} else if hasShortOptPrefix(args[o.argc]) {
				return o.getoptShort(args)
			} else if !hasOptPrefix(arg) {
				o.Value = arg
				o.argc++
				return OPT_RET_NO_ARG
			} else {
				o.Err = fmt.Errorf("Invalid option \"%s\"", arg)
				return OPT_RET_END
			}
		} else {
			return o.getoptShort(args)
		}
	}

	return OPT_RET_END
}

func (o *Options) getoptLong(args []string) int {
	fmt.Printf("getoptLong called, argc: %d, pos: %d\n", o.argc, o.pos)
	arg := args[o.argc][2:]

	tokens := strings.SplitN(arg, "=", 2)

	for i, opt := range o.opts {
		if opt.hasLong() && opt.OptLong == tokens[0] {
			//check args has an argument if the opt needs
			if opt.HasArg == OPT_REQUIRE_ARG {
				if len(tokens) == 2 {
					// --opt=value
					o.Value = tokens[1]
					o.argc++
					return i
				} else if o.argc < len(args)-1 {
					// --opt value
					o.Value = args[o.argc+1]
					o.argc += 2
					return i
				} else {
					o.Err = fmt.Errorf("Option \"--%s\" needs an argument", arg)
					return OPT_RET_END
				}
			} else {
				if len(tokens) == 2 {
					o.Err = fmt.Errorf("Invalid argument %s for opt \"--%s\"", tokens[1], tokens[0])
					return OPT_RET_END
				}
				o.argc++
				return i
			}
		}
	}

	//not found
	o.Err = fmt.Errorf("Unknown option \"--%s\"", arg)
	return OPT_RET_END
}

func (o *Options) getoptShort(args []string) int {
	arg := args[o.argc]
	tokens := strings.SplitN(arg, "=", 2)

	//ignore option prefix
	for tokens[0][o.pos] == '-' && o.pos < len(tokens[0]) {
		o.pos++
	}

	if o.pos == len(tokens[0]) {
		o.Err = fmt.Errorf("Unknown options %s", arg)
		return OPT_RET_END
	}

	for i, opt := range o.opts {
		if opt.hasShort() && opt.OptShort == tokens[0][o.pos] {
			//found it
			//check if the opt is the last one if it needs an argument
			if opt.HasArg == OPT_REQUIRE_ARG {
				if o.pos != len(tokens[0])-1 {
					o.Err = fmt.Errorf("Option \"-%c\" needs an argument", tokens[0][o.pos])
					return OPT_RET_END
				}
				if len(tokens) == 2 {
					o.Value = tokens[1]
					o.argc += 1
					o.pos = 0
					return i
				} else if o.argc != len(args)-1 {
					o.Value = args[o.argc+1]
					o.argc += 2
					o.pos = 0
					return i
				} else {
					o.Err = fmt.Errorf("Option \"-%c\" needs an argument", arg[o.pos])
					return OPT_RET_END
				}
			} else {
				if o.pos == len(tokens[0])-1 && len(tokens) == 2 {
					o.Err = fmt.Errorf("Invalid argument %s for opt \"-%c\"", tokens[1], tokens[0][o.pos])
					return OPT_RET_END
				}
				o.pos++
				if o.pos == len(tokens[0]) {
					o.pos = 0
					o.argc++
				}
				return i
			}
		}
	}

	//not found
	o.Err = fmt.Errorf("Unknown option \"-%c\"", arg[o.pos])
	return OPT_RET_END
}
