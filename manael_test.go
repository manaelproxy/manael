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
	}{
		{
			accept:      "image/*,*/*",
			path:        "/logo.png",
			statusCode:  http.StatusOK,
			contentType: "image/png",
		},
		{
			accept:      "image/webp,image/*,*/*",
			path:        "/logo.png",
			statusCode:  http.StatusOK,
			contentType: "image/webp",
		},
		{
			accept:      "image/webp,image/*,*/*",
			path:        "/empty.txt",
			statusCode:  http.StatusOK,
			contentType: "text/plain; charset=utf-8",
		},
	}

	for _, tc := range tests {
		req := httptest.NewRequest("GET", tc.path, nil)
		req.Header.Add("Accept", tc.accept)

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status Code is %d, want %d", got, want)
		}

		if got, want := resp.Header.Get("Content-Type"), tc.contentType; got != want {
			t.Errorf("Content-Type is %s, want %s", got, want)
		}
	}
}
