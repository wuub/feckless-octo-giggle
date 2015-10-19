package main

import "bytes"

func PrefixFilter(src string) (f Filter, err error) {
	prefix := []byte(src)
	return func(line []byte) bool {
		return bytes.HasPrefix(line, prefix)
	}, nil

}

func ContainsFilter(src string) (f Filter, err error) {
	part := []byte(src)
	return func(line []byte) bool {
		return bytes.Contains(line, part)
	}, nil

}
