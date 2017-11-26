package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var bind = flag.String("bind", "0.0.0.0", "")
var port = flag.String("port", "8080", "")
var upstreamURL = flag.String("upstream-url", "http://localhost:9000", "")

func main() {
	flag.Parse()

	h := &Handler{
		upstreamURL: *upstreamURL,
	}

	addr := fmt.Sprintf("%s:%s", *bind, *port)
	err := http.ListenAndServe(addr, h)
	if err != nil {
		log.Fatal(err)
	}
}
