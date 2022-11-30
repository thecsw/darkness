package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

func tools() {
	options := getEmiliaOptions(flag.NewFlagSet("tools", flag.ExitOnError))
	options.Dev = true
	emilia.InitDarkness(options)
	inputFilenames := make(chan string, customChannelCapacity)

	pages := gana.GenericWorkers(gana.GenericWorkers(inputFilenames,
		func(v string) gana.Tuple[string, string] {
			data, err := ioutil.ReadFile(filepath.Clean(v))
			if err != nil {
				fmt.Printf("Failed to open %s: %s\n", v, err.Error())
			}
			return gana.NewTuple(v, string(data))
		}, 1, customChannelCapacity), func(v gana.Tuple[string, string]) *yunyun.Page {
		return emilia.ParserBuilder.BuildParser(fdb(v.UnpackRef())).Parse()
	}, customNumWorkers, customChannelCapacity)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go emilia.FindFilesByExt(inputFilenames, workDir, wg)
	galleryDirs := map[string]bool{}
	go func(wg *sync.WaitGroup) {
		for page := range pages {
			for _, gc := range gana.Filter(func(v *yunyun.Content) bool { return v.IsGallery() }, page.Contents) {
				galleryDirs[emilia.JoinPath(filepath.Join(page.URL, gc.Summary))] = true
			}
			wg.Done()
		}
		wg.Done()
	}(wg)

	wg.Wait()

	for key := range galleryDirs {
		fmt.Println(key)
	}
}
