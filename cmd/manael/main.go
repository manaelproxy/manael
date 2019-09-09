// Copyright (c) 2017 Yamagishi Kazutoshi <ykzts@desire.sh>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main // import "manael.org/x/manael/cmd/manael"

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"manael.org/x/manael"
)

const (
	// DefaultPort is returned by default port.
	DefaultPort = 8080

	// DefaultUpstreamURL is returned by default upstream URL.
	DefaultUpstreamURL = "http://localhost:9000"
)

type config struct {
	httpAddr    string
	upstreamURL string
}

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	conf := config{}

	fs.StringVar(&conf.httpAddr, "http", "", "HTTP server address")
	fs.StringVar(&conf.upstreamURL, "upstream_url", "", "Upstream URL for processing images")

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Error: %v", err)
	}

	if conf.httpAddr == "" {
		port := os.Getenv("PORT")

		if port != "" {
			conf.httpAddr = fmt.Sprintf(":%s", port)
		} else {
			conf.httpAddr = fmt.Sprintf(":%d", DefaultPort)
		}
	}

	if conf.upstreamURL == "" {
		u := os.Getenv("MANAEL_UPSTREAM_URL")

		if u != "" {
			conf.upstreamURL = u
		} else {
			conf.upstreamURL = DefaultUpstreamURL
		}
	}

	p := manael.NewServeProxy(conf.upstreamURL)
	loggedProxy := handlers.CombinedLoggingHandler(os.Stdout, p)

	err := http.ListenAndServe(conf.httpAddr, loggedProxy)
	if err != nil {
		log.Fatal(err)
	}
}
