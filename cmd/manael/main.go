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

var bind = flag.String("bind", "0.0.0.0", "")
var port = flag.String("port", "8080", "")
var upstreamURL = flag.String("upstream-url", "http://localhost:9000", "")

func main() {
	flag.Parse()

	h := &manael.Handler{
		UpstreamURL: *upstreamURL,
	}

	addr := fmt.Sprintf("%s:%s", *bind, *port)
	err := http.ListenAndServe(addr, handlers.LoggingHandler(os.Stdout, h))
	if err != nil {
		log.Fatal(err)
	}
}
