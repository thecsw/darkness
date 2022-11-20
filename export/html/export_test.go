package html

import (
	"io/ioutil"
	"testing"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/parse/orgmode"
	_ "github.com/thecsw/darkness/parse/orgmode"
)

const (
	testFileHome = "../../parse/orgmode/testfiles/home.org"
	testFileArch = "../../parse/orgmode/testfiles/arch.org"
)

func BenchmarkAExportHome(b *testing.B) {
	emilia.InitDarkness(&emilia.EmiliaOptions{Test: true})
	data, err := ioutil.ReadFile(testFileHome)
	if err != nil {
		b.Fatal(err)
	}
	page := orgmode.ParserOrgmodeBuilder{}.
		BuildParser(testFileHome, string(data)).Parse()
	for i := 0; i < b.N; i++ {
		exporterBuilder.BuildExporter(page).Export()
	}
}

func BenchmarkAExportArch(b *testing.B) {
	emilia.InitDarkness(&emilia.EmiliaOptions{Test: true})
	data, err := ioutil.ReadFile(testFileArch)
	if err != nil {
		b.Fatal(err)
	}
	page := orgmode.ParserOrgmodeBuilder{}.
		BuildParser(testFileArch, string(data)).Parse()
	for i := 0; i < b.N; i++ {
		exporterBuilder.BuildExporter(page).Export()
	}
}
