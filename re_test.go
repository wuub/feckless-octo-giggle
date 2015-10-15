package main

import (
	"testing"
)

func TestReFilterTrue(t *testing.T) {
	f, _ := ReFilter(".*")
	res := f([]byte("hi"))
	if res != true {
		t.Fatal("result shoult be true for `true` script")
	}
}

func TestReFilterFalse(t *testing.T) {
	f, _ := ReFilter(`^nope$`)
	res := f([]byte("hi"))
	if res != false {
		t.Fatal("result shoult be false for `false` script")
	}
}

func TestReFilterEvenOdd(t *testing.T) {
	f, err := ReFilter(".{5,10}")
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

func BenchmarkReFilter(b *testing.B) {
	f, err := ReFilter(`.*`)
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
