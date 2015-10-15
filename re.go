package main

import (
	"regexp"
)

func ReFilter(src string) (f Filter, err error) {
	re, err := regexp.Compile(src)
	if err != nil {
		return
	}

	f = func(line []byte) bool {
		return re.Match(line)
	}
	return f, nil
}
