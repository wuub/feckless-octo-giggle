package main

import (
	"testing"
)

func TestPrefixFilterTrue(t *testing.T) {
	f, _ := PrefixFilter("")
	res := f([]byte("hi"))
	if res != true {
		t.Fatal("result shoult be true for `true` script")
	}
}

func TestPrefixFilterFalse(t *testing.T) {
	f, _ := PrefixFilter(`nope`)
	res := f([]byte("hi"))
	if res != false {
		t.Fatal("result shoult be false for `false` script")
	}
}

func TestPrefixFilterEvenOdd(t *testing.T) {
	f, err := PrefixFilter("abcde")
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

func BenchmarkPrefixFilter(b *testing.B) {
	f, err := PrefixFilter(`laksdkl`)
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
