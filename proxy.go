// Copyright (c) 2019 Yamagishi Kazutoshi <ykzts@desire.sh>
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

package manael // import "manael.org/x/manael"

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// A Proxy responds to an HTTP request.
type Proxy struct {
	Transport http.RoundTripper
}

func copyHeaders(w http.ResponseWriter, resp *http.Response) {
	for k, values := range resp.Header {
		if k != "Vary" {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}

	keys := []string{"Accept"}
	for _, v := range strings.Split(resp.Header.Get("Vary"), ",") {
		v = strings.TrimSpace(v)

		if !strings.EqualFold(v, "Accept") {
			keys = append(keys, v)
		}
	}

	w.Header().Set("Vary", strings.Join(keys[:], ", "))
	w.Header().Set("Server", "Manael")
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := p.Transport.RoundTrip(r)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		log.Printf("error: %v\n", err)

		return
	}
	defer resp.Body.Close()

	copyHeaders(w, resp)

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}
