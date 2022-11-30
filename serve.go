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
	port := serveCmd.Int("port", 8080, "port number to use (default 8080)")
	options := getEmiliaOptions(serveCmd)
	options.URL = "http://127.0.0.1:" + strconv.Itoa(*port)
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
