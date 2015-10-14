package main

import (
	"testing"
)

func TestJsFilterSimple(t *testing.T) {
	f, err := JsFilter(`function filter(line) { return true; }`)
	if err != nil {
		t.Fatal(err)
	}

	res := f([]byte("hi"))
	if res != true {
		t.Fatal("result shoult be true for `true` script")
	}
}

func BenchmarkJsFilter(b *testing.B) {
	f, err := JsFilter(`function filter(line) { return true; }`)
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
