// Copyright (c) 2018 Yamagishi Kazutoshi <ykzts@desire.sh>
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
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	cdnURL = "https://cdn.example.com"
)

func TestNewServeProxy(t *testing.T) {
	p := NewServeProxy(cdnURL)

	if got, want := p.UpstreamURL.String(), cdnURL; got != want {
		t.Errorf("Upstream URL is %s, want %s", got, want)
	}
}

func TestServeProxy_ServeHTTP(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/logo.png")
	})
	mux.HandleFunc("/empty.jpeg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/empty.jpeg")
	})
	mux.HandleFunc("/empty.gif", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/empty.gif")
	})
	mux.HandleFunc("/empty.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/empty.txt")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	p := NewServeProxy(ts.URL)

	tests := []struct {
		accept      string
		path        string
		statusCode  int
		contentType string
		checksum    string
	}{
		{
			accept:      "image/*,*/*",
			path:        "/logo.png",
			statusCode:  http.StatusOK,
			contentType: "image/png",
			checksum:    "87209ba2999bf451589e6333d0699dfcf6fa7c07eda5f4aa39051787dca62f2a",
		},
		{
			accept:      "image/webp,image/*,*/*",
			path:        "/logo.png",
			statusCode:  http.StatusOK,
			contentType: "image/webp",
			checksum:    "7af3f597f7965426d92f0333d27cf39377c9415107b2ffa716f72b7fe28ba2e9",
		},
		{
			accept:      "image/*,*/*",
			path:        "/empty.jpeg",
			statusCode:  http.StatusOK,
			contentType: "image/jpeg",
			checksum:    "1a8b498fcc782ef585778f6cd29640d78f0d2a25371786dcd62b61867d382b94",
		},
		{
			accept:      "image/webp,image/*,*/*",
			path:        "/empty.jpeg",
			statusCode:  http.StatusOK,
			contentType: "image/webp",
			checksum:    "69e3fa438596fdb60d6dee6aea10a92ebe2f19d50c4b069a8aaa97aa06b1a255",
		},
		{
			accept:      "image/*,*/*",
			path:        "/empty.gif",
			statusCode:  http.StatusOK,
			contentType: "image/gif",
			checksum:    "a065920df8cc4016d67c3a464be90099c9d28ffe7c9e6ee3a18f257efc58cbd7",
		},
		{
			accept:      "image/webp,image/*,*/*",
			path:        "/empty.gif",
			statusCode:  http.StatusOK,
			contentType: "image/gif",
			checksum:    "a065920df8cc4016d67c3a464be90099c9d28ffe7c9e6ee3a18f257efc58cbd7",
		},
		{
			accept:      "image/webp,image/*,*/*",
			path:        "/empty.txt",
			statusCode:  http.StatusOK,
			contentType: "text/plain; charset=utf-8",
			checksum:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}

	for _, tc := range tests {
		req := httptest.NewRequest("GET", tc.path, nil)
		req.Header.Add("Accept", tc.accept)

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status Code is %d, want %d", got, want)
		}

		if got, want := resp.Header.Get("Content-Type"), tc.contentType; got != want {
			t.Errorf("Content-Type is %s, want %s", got, want)
		}

		h := sha256.New()
		io.Copy(h, resp.Body)

		if got, want := fmt.Sprintf("%x", h.Sum(nil)), tc.checksum; got != want {
			t.Errorf("Chacksum is %s, want %s", got, want)
		}
	}
}
