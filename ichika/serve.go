package ichika

import (
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/ichika/hizuru"
	"github.com/thecsw/darkness/v3/ichika/kuroko"
	"github.com/thecsw/darkness/v3/yunyun"
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
	options.Url = "http://127.0.0.1:" + strconv.Itoa(*port)
	// Override the output extension to .html
	options.OutputExtension = puck.ExtensionHtml
	// emilia.InitDarkness(options)
	conf := alpha.BuildConfig(options)

	puck.Logger.SetPrefix("Server 🍩 ")

	build(conf)
	// disable akane after the first build
	kuroko.Akaneless = true
	puck.Logger.Print("Serving the files", "url", options.Url)

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
	fileServer(r, "/", http.Dir(string(conf.Runtime.WorkDir)))

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
		// Validate URL before passing it to exec.Command
		if isURLSafe(options.Url) {
			if err := exec.Command("open", options.Url).Run(); err != nil {
				puck.Logger.Error("Couldn't open the browser", err)
			}
		} else {
			puck.Logger.Error("Invalid URL for browser", "url", options.Url)
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
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Rename) ||
					event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
					puck.Logger.Warn("A file was modified", "path", filename)
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
	for _, inputFilenameToWatch := range hizuru.FindFilesByExtSimple(conf) {
		err = watcher.Add(string(inputFilenameToWatch))
		if err != nil {
			log.Fatal(err)
		}
	}
	puck.Logger.Print("Listening to file changes", "num", len(watcher.WatchList()), "dir", conf.Runtime.WorkDir)

	puck.Logger.Print("Press Ctrl-C to stop the server")
	// Block main goroutine forever.
	<-make(chan struct{})
}

// isURLSafe checks if a URL is safe to pass to exec.Command
// It verifies the URL is properly formatted and is a localhost URL.
func isURLSafe(urlString string) bool {
	// Only allow http and https URLs
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		return false
	}

	// Parse the URL to validate its format
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return false
	}

	// Only allow localhost URLs
	host := parsedURL.Hostname()
	if host != "localhost" && host != "127.0.0.1" {
		return false
	}

	return true
}

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// Taken from https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("fileServer does not permit any Url parameters.")
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
