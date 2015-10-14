package main

import (
	"github.com/yuin/gopher-lua"
)

func LuaFilter(src string) (f Filter, err error) {
	l := lua.NewState()
	l.OpenLibs()
	err = l.DoString(src)
	if err != nil {
		return
	}
	fn := l.GetGlobal("filter")

	f = func(line []byte) bool {
		if err := l.CallByParam(lua.P{
			Fn:      fn,
			NRet:    1,
			Protect: false,
		}, lua.LString(line)); err != nil {
			panic(err)
		}
		ret := l.Get(-1)
		l.Pop(1)
		return lua.LVAsBool(ret)
	}
	return f, nil
}
