package ichika

import (
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
	"github.com/thecsw/darkness/emilia/alpha"
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
	port := serveCmd.Int("port", defaultServePort, "port number to use")
	noBrowser := serveCmd.Bool("no-browser", false, "do not open the browser")
	options := getAlphaOptions(serveCmd)
	options.URL = "http://127.0.0.1:" + strconv.Itoa(*port)
	// Override the output extension to .html
	options.OutputExtension = puck.ExtensionHtml
	//emilia.InitDarkness(options)
	conf := alpha.BuildConfig(options)

	puck.Logger.SetPrefix("Server 🍩 ")

	build(conf)
	puck.Logger.Print("Serving the files", "url", options.URL)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.NoCache)
	r.Use(middleware.Recoverer)

	// Set up the server's timeouts.
	srv := &http.Server{
		Addr:              "127.0.0.1:" + strconv.Itoa(*port),
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Tune it to serve local files.
	FileServer(r, "/", http.Dir(string(conf.Runtime.WorkDir)))

	// Spin the local server up.
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	// File watcher will rebuild dir if any files change.
	go launchWatcher(conf)
	puck.Logger.Print("Launched file watcher")

	// Try to open the local server with `open` command.
	if !*noBrowser {
		time.Sleep(500 * time.Millisecond)
		if err := exec.Command("open", options.URL).Run(); err != nil {
			puck.Logger.Error("Couldn't open the browser", err)
		}
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	puck.Logger.Print("Shutting down the server + cleaning up")
	isQuietMegumin = true
	removeOutputFiles(conf)
	puck.Logger.Print("farewell")
}

// launchWatcher watches for any file creations, changes, modifications, deletions
// and rebuilds the directory as that happens.
func launchWatcher(conf *alpha.DarknessConfig) {
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
					puck.Logger.Warn("stopped watching")
					return
				}
				filename := string(conf.Runtime.WorkDir.Rel(yunyun.FullPathFile(event.Name)))
				if strings.HasSuffix(filename, conf.Project.Output) ||
					strings.HasPrefix(filepath.Base(filename), `.`) {
					continue
				}
				// Skip CHMOD events that IDE and editors do by default
				if event.Has(fsnotify.Chmod) {
					continue
				}
				if event.Has(fsnotify.Write) {
					puck.Logger.Warn("A file was modified", "path", filename)
				}
				if event.Has(fsnotify.Create) {
					puck.Logger.Warn("A file was created", "path", filename)
				}
				if event.Has(fsnotify.Remove) {
					puck.Logger.Warn("A file was removed", "path", filename)
				}
				if event.Has(fsnotify.Rename) {
					puck.Logger.Warn("A file was renamed", "path", filename)
				}
				build(conf)
			case err, ok := <-watcher.Errors:
				if !ok {
					puck.Logger.Warn("Watcher is leaving")
					return
				}
				puck.Logger.Error("Watcher", "err", err)
			}
		}
	}()

	// start adding all the source files
	for _, inputFilenameToWatch := range FindFilesByExtSimple(conf) {
		err = watcher.Add(string(inputFilenameToWatch))
		if err != nil {
			log.Fatal(err)
		}
	}
	puck.Logger.Print("Listening to file changes", "num", len(watcher.WatchList()), "dir", workDir)

	puck.Logger.Print("Press Ctrl-C to stop the server")
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
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
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
