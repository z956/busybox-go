package main

import (
	"fmt"
	"strconv"
	"time"
)

type sleep struct {
	dur uint
}

func (s *sleep) parseTimeOpt(dur string) error {
	last := dur[len(dur)-1]
	mul := 1.0
	switch last {
	case 's':
		dur = dur[:len(dur)-1]
	case 'm':
		dur = dur[:len(dur)-1]
		mul = 60
	case 'h':
		dur = dur[:len(dur)-1]
		mul = 60 * 60
	case 'd':
		dur = dur[:len(dur)-1]
		mul = 60 * 60 * 24
	default:
	}

	size, err := strconv.ParseFloat(dur, 64)
	s.dur = (uint)(size * mul)
	return err
}

func (s *sleep) run() int {
	time.Sleep(time.Duration(s.dur) * time.Second)
	return 0
}

func usage() {
	fmt.Println("Usage: sleep [N]")
	fmt.Println()

	fmt.Println("Pause for a time equal to the total of the args given, where each arg can")
	fmt.Println("have an optional suffix of (s)econds, (m)inutes, (h)ours, or (d)ays")
}

func SleepMain(args []string) int {
	if len(args) == 1 {
		usage()
		return -1
	}

	s := sleep{}
	if err := s.parseTimeOpt(args[1]); err != nil {
		fmt.Println(err)
		fmt.Println()
		usage()
		return -1
	} else {
		return s.run()
	}
}
