package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

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
	serveCmd := flag.NewFlagSet(serveCommand, flag.ExitOnError)
	port := serveCmd.Int("port", defaultServePort, "port number to use (default 8080)")
	options := getEmiliaOptions(serveCmd)
	options.URL = "http://127.0.0.1:" + strconv.Itoa(*port)
	// Override the output extension to .html
	options.OutputExtension = puck.ExtensionHtml
	emilia.InitDarkness(options)
	start := time.Now()
	build()
	fmt.Printf("Built in %d ms\n\n", time.Since(start).Milliseconds())
	fmt.Println("Serving on", options.URL)
	go func() {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), http.FileServer(http.Dir(workDir))))
	}()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	fmt.Println("Shutting down the server + cleaning up")
	isQuietMegumin = true
	removeOutputFiles()
	fmt.Println("farewell")
}
