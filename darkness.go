package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/sanity-io/litter"
)

const (
	workDir = "sandyuraz"
)

var (
	conf DarknessConfig
)

func main() {

	confData, _ := ioutil.ReadFile("darkness.toml")
	_, err := toml.Decode(string(confData), &conf)
	if err != nil {
		panic(err)
	}

	orgfiles, err := findOrgFiles(workDir)
	if err != nil {
		panic(err)
	}
	litter.Dump(orgfiles)

	for _, file := range orgfiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		lines := strings.Split(string(data), "\n")
		page := Parse(lines)
		page.URL = conf.URL + strings.TrimPrefix(filepath.Dir(file), workDir) + "/"
		fmt.Println("Processed", file, ":", page.URL)
		ioutil.WriteFile(filepath.Dir(file)+"/index.html", []byte(buildHTML(page)), 0644)
	}
}

func findOrgFiles(dir string) ([]string, error) {
	files := make([]string, 0, 32)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".org" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
