package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

// misaCommandFunc will support many different tools that darkness can support,
// such as creating gallery previews, etc. WIP.
func misaCommandFunc() {
	misaCmd := flag.NewFlagSet(misaCommand, flag.ExitOnError)

	buildGalleryPreviews := misaCmd.Bool("gallery-previews", false, "build gallery previews")
	removeGalleryPreviews := misaCmd.Bool("no-gallery-previews", false, "delete gallery previews")

	options := getEmiliaOptions(misaCmd)
	options.Dev = true
	emilia.InitDarkness(options)

	if *buildGalleryPreviews {
		buildGalleryFiles()
		return
	}
	if *removeGalleryPreviews {
		removeGalleryFiles()
		return
	}
	fmt.Println("I don't know what you want me to do, see -help")
}

const (
	galleryPreviewImageSize = 500
	galleryPreviewImageBlur = 20
)

// buildGalleryFiles finds all the gallery entries and build a resized
// preview version of it.
func buildGalleryFiles() {
	for _, galleryFile := range getGalleryFiles() {
		img, err := imaging.Open(galleryFile, imaging.AutoOrientation(false))
		if err != nil {
			fmt.Println("failed to open", galleryFile, ":", err.Error())
			continue
		}

		// Resize the image to save up on storage.
		img = imaging.Resize(img, galleryPreviewImageSize, 0, imaging.Lanczos)

		blurred := imaging.Blur(img, galleryPreviewImageBlur)
		newFile := emilia.GalleryPreview(galleryFile)
		fmt.Println("Saving", newFile)
		err = imaging.Save(blurred, newFile)
		if err != nil {
			fmt.Println("failed to save", newFile)
			continue
		}
	}
}

func removeGalleryFiles() {
	for _, galleryFile := range getGalleryFiles() {
		newFile := emilia.GalleryPreview(galleryFile)
		if err := os.Remove(newFile); err != nil && !os.IsNotExist(err) {
			fmt.Println("Couldn't delete", newFile, "| reason:", err.Error())
		}
	}
}

func getGalleryFiles() []string {
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
	// Launch a second discovery for gallery files
	galleryFiles := make([]string, 0, 32)
	go func(wg *sync.WaitGroup) {
		for page := range pages {
			for _, gc := range page.Contents.Galleries() {
				// If it's external, can't delete it, so skip
				if gc.GalleryPathIsExternal {
					continue
				}
				for _, item := range gc.List {
					_, link, _, _ := yunyun.ExtractLink(item)
					galleryFiles = append(galleryFiles, emilia.JoinPath(
						filepath.Join(page.URL, gc.GalleryPath, link),
					))
				}
			}
			wg.Done()
		}
		wg.Done()
	}(wg)

	wg.Wait()
	return galleryFiles
}
