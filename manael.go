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

// Package manael provides HTTP handler for processing images.
package manael // import "manael.org/x/manael"

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// A ServeProxy responds to an HTTP request.
type ServeProxy struct {
	UpstreamURL *url.URL
}

func copyHeaders(w http.ResponseWriter, resp *http.Response) {
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}

func (p *ServeProxy) transfer(w http.ResponseWriter, resp *http.Response) {
	copyHeaders(w, resp)

	w.Header().Set("Vary", "Accept")
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func shouldEncodeToWebP(resp *http.Response) bool {
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified) {
		return false
	}

	contentType := resp.Header.Get("Content-Type")
	return contentType == "image/jpeg" || contentType == "image/png"
}

func canDecodeWebP(r *http.Request) bool {
	accepts := r.Header.Get("Accept")

	for _, v := range strings.Split(accepts, ",") {
		t := strings.TrimSpace(v)
		if strings.HasPrefix(t, "image/webp") {
			return true
		}
	}

	return false
}

func (p *ServeProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s%s", p.UpstreamURL.String(), r.URL.RequestURI())

	resp, err := request(url, r)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		log.Println(err)

		return
	}
	defer resp.Body.Close()

	if !(shouldEncodeToWebP(resp) && canDecodeWebP(r)) {
		p.transfer(w, resp)

		return
	}

	if resp.StatusCode == http.StatusNotModified {
		copyHeaders(w, resp)

		w.Header().Set("Cache-Control", resp.Header.Get("Cache-Control"))
		w.WriteHeader(http.StatusNotModified)

		return
	}

	buf, err := convert(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)

		return
	}

	copyHeaders(w, resp)

	w.Header().Set("Content-Type", "image/webp")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.Header().Set("Vary", "Accept")
	w.WriteHeader(http.StatusOK)

	io.Copy(w, buf)
}

// NewServeProxy returns a new ServeProxy given a upstream URL
func NewServeProxy(rawURL string) *ServeProxy {
	u, _ := url.Parse(rawURL)

	return &ServeProxy{u}
}
