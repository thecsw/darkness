package orgmode

import (
	"io/ioutil"
	"testing"

	"github.com/thecsw/darkness/emilia"
)

const (
	testFileHome = "./testfiles/home.org"
	testFileArch = "./testfiles/arch.org"
)

func BenchmarkAParseHome(b *testing.B) {
	emilia.InitDarkness(&emilia.EmiliaOptions{Test: true})
	data, err := ioutil.ReadFile(testFileHome)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		ParserOrgmode{}.Parse(testFileHome, string(data))
	}
}

func BenchmarkBParseArch(b *testing.B) {
	emilia.InitDarkness(&emilia.EmiliaOptions{Test: true})
	data, err := ioutil.ReadFile(testFileArch)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		ParserOrgmode{}.Parse(testFileArch, string(data))
	}
}
