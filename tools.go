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

// toolsCommandFunc will support many different tools that darkness can support,
// such as creating gallery previews, etc. WIP.
func toolsCommandFunc() {
	options := getEmiliaOptions(flag.NewFlagSet(toolsCommand, flag.ExitOnError))
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
	go emilia.FindFilesByExt(inputFilenames, workDir, emilia.Config.Project.Input, wg)
	galleryDirs := map[string]bool{}
	go func(wg *sync.WaitGroup) {
		for page := range pages {
			for _, gc := range page.Contents.Galleries() {
				for _, item := range gc.List {
					_, link, _, _ := yunyun.ExtractLink(item)
					fmt.Println(emilia.JoinPath(filepath.Join(gc.Summary, link)))
				}
			}
			wg.Done()
		}
		wg.Done()
	}(wg)

	wg.Wait()

	// Launch a second discovery for gallery files
	galleryFiles := make(chan string, customChannelCapacity)

	for key := range galleryDirs {
		go emilia.FindFilesByExt(galleryFiles, key, ".png", wg)
	}

	for galleryFile := range galleryFiles {
		fmt.Println(galleryFile)
		wg.Done()
	}
	wg.Wait()
}
