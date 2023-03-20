package ichika

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
	"github.com/thecsw/komi"
)

const (
	// savePerms tells us what permissions to use for the
	// final export files.
	savePerms = fs.FileMode(0644)
)

var (
	// vendorGalleryImages is a flag that dictates whether we should
	// store a local copy of all remote gallery images and stub them
	// in the gallery links instead of the remote links.
	//
	// Turning this option on would result in a VERY slow build the
	// first time, as it would need to retrieve however many images
	// from remote services.
	//
	// All images will be put in "darkness_vendor" directory, which
	// will be skipped in discovery process AND should be put it
	// .gitignore by user, so they don't pollute their git objects.
	vendorGalleryImages bool
)

// OneFileCommandFunc builds a single file.
func OneFileCommandFunc() {
	fileCmd := darknessFlagset(oneFileCommand)
	fileCmd.StringVar(&filename, "input", "index.org", "file on input")
	emilia.InitDarkness(getEmiliaOptions(fileCmd))
	fmt.Println(emilia.InputToOutput(emilia.JoinWorkdir(yunyun.RelativePathFile(filename))))
}

// build builds the entire directory.
func BuildCommandFunc() {
	cmd := darknessFlagset(buildCommand)
	emilia.InitDarkness(getEmiliaOptions(cmd))
	build()
	fmt.Println("farewell")
}

// build uses set flags and emilia data to build the local directory.
func build() {

	// Create the pool that reads files and returns their handles.
	filesPool := komi.NewWithSettings(komi.Work(openPage), &komi.Settings{
		Name:     "Reading 📚 ",
		Laborers: customNumWorkers,
		Debug:    debugEnabled,
	})

	// Create a pool that take a files handle and parses it out into yunyun pages.
	parserPool := komi.NewWithSettings(komi.Work(parsePage), &komi.Settings{
		Name:     "Parsing 🧹 ",
		Laborers: customNumWorkers,
		Debug:    debugEnabled,
	})

	// Create a pool that that takes yunyun pages and exports them into request format.
	exporterPool := komi.NewWithSettings(komi.Work(exportPage), &komi.Settings{
		Name:     "Exporting 🥂 ",
		Laborers: customNumWorkers,
		Debug:    debugEnabled,
	})

	// Create a pool that reads the exported data and writes them to target files.
	writerPool := komi.NewWithSettings(komi.WorkSimpleWithErrors(writePage), &komi.Settings{
		Name:     "Writing 🎸",
		Laborers: 1,
		Debug:    debugEnabled,
	})

	// Handle any errors from the writing.
	go func() {
		for err := range writerPool.Errors() {
			fmt.Printf("failed to write %s: %s", err.Job.First, err.Error)
		}
	}()

	// Connect all the pools between each other, so the relationship is as follows,
	//
	//           Reading 📚                      Parsing 🧹
	//   path  ┌───────────┐   file handler   ┌─────────────┐
	// ──────> │ filesPool │ ───────────────> │  parserPool │
	//	   └───────────┘	          └─────────────┘
	//	    log errors				│
	//						│    parsed files
	//					       	│  aka yunyun pages
	//                                              │
	//   file  ┌────────────┐  exported data  ┌──────────────┐
	//  <───── │ writerPool │ <────────────── │ exporterPool │
	// 	   └────────────┘              	  └──────────────┘
	//           Writing 🎸                     Exporting 🥂
	//
	filesPool.Connect(parserPool)
	parserPool.Connect(exporterPool)
	exporterPool.Connect(writerPool)

	start := time.Now()

	<-emilia.FindFilesByExt(filesPool, emilia.Config.Project.Input)

	writerPool.Close()

	finish := time.Now()

	// Clear the download progress bar if present by wiping out the line.
	fmt.Print("\r\033[2K")

	fmt.Printf("Processed %d files in %d ms\n", exporterPool.JobsCompleted(), finish.Sub(start).Milliseconds())
}

//go:inline
func openPage(v yunyun.FullPathFile) gana.Tuple[yunyun.FullPathFile, *os.File] {
	return openFile(v)
}

//go:inline
func parsePage(v gana.Tuple[yunyun.FullPathFile, *os.File]) *yunyun.Page {
	return emilia.ParserBuilder.BuildParserReader(emilia.FullPathToWorkDirRel(v.First), v.Second).Parse()
}

//go:inline
func exportPage(v *yunyun.Page) gana.Tuple[string, *bufio.Reader] {
	return gana.NewTuple(emilia.InputFilenameToOutput(emilia.JoinWorkdir(v.File)), emilia.EnrichExportPageAsBufio(v))
}

//go:inline
func writePage(v gana.Tuple[string, *bufio.Reader]) error {
	_, err := writeFile(v.First, v.Second)
	return err
}
