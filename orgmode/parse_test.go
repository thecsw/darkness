package orgmode

import (
	"io/ioutil"
	"testing"
)

func BenchmarkAParseHome(b *testing.B) {
	data, err := ioutil.ReadFile("../test/home.org")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		Parse(string(data))
	}
}

func BenchmarkBParseArch(b *testing.B) {
	data, err := ioutil.ReadFile("../test/arch.org")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		Parse(string(data))
	}
}
