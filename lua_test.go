package main

import (
	"testing"
)

func TestLuaFilterTrue(t *testing.T) {
	f, _ := LuaFilter(`function filter (line) return true end`)
	res := f([]byte("hi"))
	if res != true {
		t.Fatal("result shoult be true for `true` script")
	}
}

func TestLuaFilterFalse(t *testing.T) {
	f, _ := LuaFilter(`function filter (line) return false end`)
	res := f([]byte("hi"))
	if res != false {
		t.Fatal("result shoult be false for `false` script")
	}
}

func TestLuaFilterEvenOdd(t *testing.T) {
	f, err := LuaFilter(`function filter (line) return string.len(line) > 5 end`)
	if err != nil {
		t.Fatal(err)
	}

	if f([]byte("abcd")) != false {
		t.Fatal("should be false for short string")
	}

	if f([]byte("abcdef")) != true {
		t.Fatal("should be true for long string")
	}
}

func BenchmarkLuaFilter(b *testing.B) {
	f, err := LuaFilter(`function filter (line) return true end`)
	if err != nil {
		b.Fatal(err)
	}

	line := []byte("laksdkl alksjd laksjd laksdj lkfj slkdjf lskdfjg lskdjfg lksdjf glkjsdf lksjdfj")
	b.SetBytes(int64(len(line)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(line)
	}
}
