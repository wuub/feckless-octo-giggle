package main

import (
	"github.com/robertkrimen/otto"
)

func JsFilter(src string) (f Filter, err error) {
	vm, _, err := otto.Run(src)
	if err != nil {
		return nil, err
	}

	script, err := vm.Compile("", `filter(line)`)
	if err != nil {
		return nil, err
	}

	f = func(line []byte) bool {
		vm.Set("line", line)
		value, err := vm.Run(script)
		if err != nil {
			panic(err)
		}
		res, err := value.ToBoolean()
		if err != nil {
			panic(err)
		}
		return res
	}

	return f, nil
}
