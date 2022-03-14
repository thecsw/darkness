package emilia

import (
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	Config *DarknessConfig
)

func InitDarkness(file string) {
	Config = &DarknessConfig{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	_, err = toml.Decode(string(data), Config)
	if err != nil {
		panic(err)
	}
	// URL must end with a trailing forward slash
	if !strings.HasSuffix(Config.URL, "/") {
		Config.URL += "/"
	}
	Config.URLPath, err = url.Parse(Config.URL)
	if err != nil {
		panic(err)
	}
}

func JoinPath(elem string) string {
	u, _ := url.Parse(elem)
	return Config.URLPath.ResolveReference(u).String()
}
