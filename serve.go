package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/emilia/puck"
)

const (
	// defaultServePort is the default port used when serving
	// local files.
	defaultServePort = 8080
)

// serveCommandFunc builds the website, local serves it on 8080 and then
// cleans the files.
func serveCommandFunc() {
	serveCmd := darknessFlagset(serveCommand)
	port := serveCmd.Int("port", defaultServePort, "port number to use (default 8080)")
	options := getEmiliaOptions(serveCmd)
	options.URL = "http://127.0.0.1:" + strconv.Itoa(*port)
	// Override the output extension to .html
	options.OutputExtension = puck.ExtensionHtml
	emilia.InitDarkness(options)
	build()
	log.Println("Serving on", options.URL)
	go func() {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), http.FileServer(http.Dir(workDir))))
	}()
	go launchWatcher()
	log.Println("Launched file watcher")
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	log.Println("Shutting down the server + cleaning up")
	isQuietMegumin = true
	removeOutputFiles()
	log.Println("farewell")
}

// launchWatcher watches for any file creations, changes, modifications, deletions
// and rebuilds the directory as that happens.
func launchWatcher() {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				filename := relPathToWorkdir(event.Name)
				if strings.HasSuffix(filename, emilia.Config.Project.Output) ||
					strings.HasPrefix(filepath.Base(filename), `.`) {
					continue
				}
				// Skip CHMOD events that IDE and editors do by default
				if event.Has(fsnotify.Chmod) {
					continue
				}
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
				if event.Has(fsnotify.Create) {
					log.Println("created file:", event.Name)
				}
				if event.Has(fsnotify.Remove) {
					log.Println("removed file:", event.Name)
				}
				if event.Has(fsnotify.Rename) {
					log.Println("renamed file:", event.Name)
				}
				build()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watcher error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(workDir)
	if err != nil {
		log.Fatal(err)
	}
	// Block main goroutine forever.
	<-make(chan struct{})
}
