package main

import ()

type Readliner interface {
	ReadLine() (line []byte, isPrefix bool, err error)
}

type Filter func([]byte) bool

type Sender interface {
	Send([]byte) error
}

func Pipe(r Readliner, s Sender, f Filter) error {
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		if f(line) {
			if err = s.Send(line); err != nil {
				return err
			}
		}
	}
}
