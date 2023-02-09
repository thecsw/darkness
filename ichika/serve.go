package ichika

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/yunyun"
)

const (
	// defaultServePort is the default port used when serving
	// local files.
	defaultServePort = 8080
)

// ServeCommandFunc builds the website, local serves it on 8080 and then
// cleans the files.
func ServeCommandFunc() {
	serveCmd := darknessFlagset(serveCommand)
	port := serveCmd.Int("port", defaultServePort, "port number to use (default 8080)")
	noBrowser := serveCmd.Bool("no-browser", false, "do not open the browser")
	options := getEmiliaOptions(serveCmd)
	options.URL = "http://127.0.0.1:" + strconv.Itoa(*port)
	// Override the output extension to .html
	options.OutputExtension = puck.ExtensionHtml
	emilia.InitDarkness(options)
	build()
	log.Println("Serving on", options.URL)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.NoCache)
	r.Use(middleware.Recoverer)
	FileServer(r, "/", http.Dir(emilia.Config.WorkDir))

	go func() {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), r))
	}()
	go launchWatcher()
	log.Println("Launched file watcher")

	// Try to open the local server with `open` command.
	if !*noBrowser {
		time.Sleep(500 * time.Millisecond)
		if err := exec.Command("open", options.URL).Run(); err != nil {
			log.Println("couldn't open the browser:", err.Error())
		}
	}

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
				// fmt.Println(event)
				if !ok {
					log.Println("stopped watching")
					return
				}
				filename := string(emilia.FullPathToWorkDirRel(yunyun.FullPathFile(event.Name)))
				if strings.HasSuffix(filename, emilia.Config.Project.Output) ||
					strings.HasPrefix(filepath.Base(filename), `.`) {
					continue
				}
				// Skip CHMOD events that IDE and editors do by default
				if event.Has(fsnotify.Chmod) {
					continue
				}
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", filename)
				}
				if event.Has(fsnotify.Create) {
					log.Println("created file:", filename)
				}
				if event.Has(fsnotify.Remove) {
					log.Println("removed file:", filename)
				}
				if event.Has(fsnotify.Rename) {
					log.Println("renamed file:", filename)
				}
				build()
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Println("finished watching")
					return
				}
				fmt.Println("watcher error:", err)
			}
		}
	}()

	// start adding all the source files
	for _, toWatch := range emilia.FindFilesByExtSimple(emilia.Config.Project.Input) {
		err = watcher.Add(string(toWatch))
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Listening to %d files in %s\n\n", len(watcher.WatchList()), workDir)

	fmt.Println("Press Ctrl-C to stop the server")
	// Block main goroutine forever.
	<-make(chan struct{})
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// Taken from https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
