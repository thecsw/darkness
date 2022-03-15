package emilia

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
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
		fmt.Printf("failed to open the config %s: %s", file, err.Error())
		os.Exit(1)
	}
	_, err = toml.Decode(string(data), Config)
	if err != nil {
		fmt.Printf("failed to decode the config %s: %s", file, err.Error())
		os.Exit(1)
	}
	// If the URL is empty, then plug in the current directory
	if len(Config.URL) < 1 {
		Config.URL, err = os.Getwd()
		if err != nil {
			fmt.Printf("failed to get current directory because config url was not given: %s", err.Error())
			os.Exit(1)
		}
	}
	// URL must end with a trailing forward slash
	if !strings.HasSuffix(Config.URL, "/") {
		Config.URL += "/"
	}
	Config.URLPath, err = url.Parse(Config.URL)
	if err != nil {
		fmt.Printf("failed to parse url from config %s: %s", Config.URL, err.Error())
		os.Exit(1)
	}
}

func JoinPath(elem string) string {
	u, _ := url.Parse(elem)
	return Config.URLPath.ResolveReference(u).String()
}
