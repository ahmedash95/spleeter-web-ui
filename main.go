package main

import (
	"flag"
	"fmt"
	"net/http"
)

var port string

func main() {
	flag.StringVar(&port, "port", "8000", "service bind port")
	flag.Parse()

	validateRequirements()
	// err := split("/Users/ahmedashraf/ash-project/deezer-spliter/bh.mp3", "/Users/ahmedashraf/ash-project/deezer-spliter/output")
	// if err != nil {
	// 	panic(err)
	// }
	startWebServer()
}

func startWebServer() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("media"))))
	http.HandleFunc("/upload", mediaUpload)

	fmt.Printf("Starting server on 0.0.0.0:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
