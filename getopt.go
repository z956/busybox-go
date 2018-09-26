package main

import (
	"fmt"
	"strings"
)

const OPT_NO_ARG = 0
const OPT_REQUIRE_ARG = 1

type Option struct {
	OptLong  string
	OptShort byte
	HasArg   int
}

func (o Option) hasLong() bool {
	return o.OptLong != ""
}
func (o Option) hasShort() bool {
	return o.OptShort != 0
}

type ParseResultFunc func(opt Option, val string, userdata interface{}) error

func OptParse(args []string, opts []Option, cb ParseResultFunc, userdata interface{}) error {
	if err := checkOpts(opts); err != nil {
		return err
	}

	emptyOpt := Option{"", 0, OPT_NO_ARG}
	arglen := len(args)
	for i := 0; i < arglen; i++ {
		ss := strings.SplitN(args[i], "=", 2)

		if r, err := parseArg(opts, ss[0]); err != nil {
			return err
		} else {
			if len(r) == 0 {
				//not an option
				if err := cb(emptyOpt, args[i], userdata); err != nil {
					return err
				}
			} else {
				for _, o := range r {
					if o.HasArg == OPT_NO_ARG {
						if err := cb(o, "", userdata); err != nil {
							return err
						}
					} else if o.HasArg == OPT_REQUIRE_ARG {
						if len(ss) == 2 {
							if err := cb(o, ss[1], userdata); err != nil {
								return err
							}
						} else if i == arglen-1 {
							return fmt.Errorf("Missing argument with option '%s'", ss[0])
						} else {
							i += 1
							if err := cb(o, args[i], userdata); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func checkOpts(opts []Option) error {
	for i, o := range opts {
		if !o.hasLong() && !o.hasShort() {
			return fmt.Errorf("Empty option")
		}
		if o.hasShort() && !isOptValid(o.OptShort) {
			return fmt.Errorf("Invalid short opt: %x", o.OptShort)
		}
		if o.HasArg != OPT_NO_ARG && o.HasArg != OPT_REQUIRE_ARG {
			return fmt.Errorf("Invalid opt flag: %d", o.HasArg)
		}

		for _, oo := range opts[0:i] {
			if oo.hasLong() && o.hasLong() && oo.OptLong == o.OptLong {
				return fmt.Errorf("Duplicate long opt: %s", o.OptLong)
			}
			if oo.hasShort() && o.hasShort() && oo.OptShort == o.OptShort {
				return fmt.Errorf("Duplicate short opt: %c", o.OptShort)
			}
		}
	}
	return nil
}

func parseArg(opts []Option, arg string) ([]Option, error) {
	//arg should be strip with '='
	result := make([]Option, 0)
	if strings.HasPrefix(arg, "--") && arg != "--" {
		//long opt
		arg = arg[2:]
		if o, err := getOptionByLong(opts, arg); err != nil {
			return nil, err
		} else {
			result = append(result, o)
		}
	} else if strings.HasPrefix(arg, "-") && arg != "-" {
		//short opts, may be aggregated
		//only the last opt can has arg
		arg = arg[1:]
		arglen := len(arg)

		for i := 0; i < arglen; i++ {
			if o, err := getOptionByShort(opts, arg[i]); err != nil {
				return nil, err
			} else {
				if i != arglen-1 && o.HasArg == OPT_REQUIRE_ARG {
					return nil, fmt.Errorf("Invalid argument with opt '%c'", arg[i])
				}
				result = append(result, o)
			}
		}
	}

	return result, nil
}

func getOptionByLong(opts []Option, opt string) (Option, error) {
	for _, o := range opts {
		if o.hasLong() && opt == o.OptLong {
			return o, nil
		}
	}
	return Option{}, fmt.Errorf("Unknown option: %s", opt)
}

func getOptionByShort(opts []Option, opt byte) (Option, error) {
	for _, o := range opts {
		if o.hasShort() && opt == o.OptShort {
			return o, nil
		}
	}
	return Option{}, fmt.Errorf("Unknown option: %c", opt)
}

func isOptValid(b byte) bool {
	switch {
	case b >= 'a' && b <= 'z':
		return true
	case b >= 'A' && b <= 'Z':
		return true
	case b >= '0' && b <= '9':
		return true
	default:
		return false
	}
}
