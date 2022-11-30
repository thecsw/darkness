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
)

func serve() {
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	serveCmd.StringVar(&workDir, "dir", ".", "where do I look for files")
	serveCmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	serveCmd.BoolVar(&disableParallel, "disable-parallel", false, "disable parallel build (only use one worker)")
	serveCmd.IntVar(&customNumWorkers, "workers", defaultNumOfWorkers, "number of workers to spin up")
	serveCmd.IntVar(&customChannelCapacity, "capacity", defaultNumOfWorkers, "worker channels' capacity")
	port := serveCmd.Int("port", 8080, "port number to use (default 8080)")
	portStr := ":" + strconv.Itoa(*port)

	useCurrentDirectory = true

	if err := serveCmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("failed to parse build arguments, fatal: %s", err.Error())
		os.Exit(1)
	}

	// Read the config and initialize emilia settings.
	options := &emilia.EmiliaOptions{
		DarknessConfig: darknessToml,
		Dev:            useCurrentDirectory,
		URL:            "http://127.0.0.1" + portStr,
	}
	emilia.InitDarkness(options)

	start := time.Now()
	build()
	fmt.Printf("Built in %d ms\n\n", time.Since(start).Milliseconds())

	fmt.Println("Serving on", options.URL)
	go func() {
		log.Fatal(http.ListenAndServe(portStr, http.FileServer(http.Dir(workDir))))
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	// Block until received
	<-sigint
	fmt.Println("Shutting down the server + cleaning up")

	isQuietMegumin = true
	removeOutputFiles()
	fmt.Println("farewell")
}
