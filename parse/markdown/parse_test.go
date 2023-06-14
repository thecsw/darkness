package markdown

import (
	"os"
	"testing"

	"github.com/thecsw/darkness/emilia"
)

const (
	testFileHome = "./testfiles/home.org"
	testFileArch = "./testfiles/arch.org"
)

func BenchmarkAParseHome(b *testing.B) {
	emilia.InitDarkness(&emilia.EmiliaOptions{Test: true})
	data, err := os.ReadFile(testFileHome)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		parserBuilder.BuildParser(testFileHome, string(data)).Parse()
	}
}

func BenchmarkBParseArch(b *testing.B) {
	emilia.InitDarkness(&emilia.EmiliaOptions{Test: true})
	data, err := os.ReadFile(testFileArch)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		parserBuilder.BuildParser(testFileArch, string(data)).Parse()
	}
}
