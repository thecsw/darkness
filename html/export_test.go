package html

import (
	"io/ioutil"
	"testing"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/orgmode"
)

func BenchmarkAExportHome(b *testing.B) {
	emilia.InitDarkness("../test/darkness.toml")
	data, err := ioutil.ReadFile("../test/home.org")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		ExportPage(orgmode.Parse(string(data)))
	}
}

func BenchmarkAExportArch(b *testing.B) {
	emilia.InitDarkness("../test/darkness.toml")
	data, err := ioutil.ReadFile("../test/arch.org")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		ExportPage(orgmode.Parse(string(data)))
	}
}
