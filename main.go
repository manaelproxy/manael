package main

import (
	"flag"
	"log"
	"net/http"
)

var port = flag.String("port", "8080", "")
var upstreamURL = flag.String("upstream-url", "http://localhost:9000", "")

func main() {
	flag.Parse()

	h := &Handler{
		upstreamURL: *upstreamURL,
	}

	err := http.ListenAndServe(":"+*port, h)
	if err != nil {
		log.Fatal(err)
	}
}
