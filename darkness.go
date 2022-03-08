package main

import (
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/sanity-io/litter"
)

var conf DarknessConfig

func main() {
	data, err := ioutil.ReadFile("flcl.org")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")
	page := Parse(lines)

	litter.Dump(page)

	confData, _ := ioutil.ReadFile("darkness.toml")
	_, err = toml.Decode(string(confData), &conf)
	if err != nil {
		panic(err)
	}
	//litter.Dump(conf)

	//fmt.Println(buildHTML(page))
}
