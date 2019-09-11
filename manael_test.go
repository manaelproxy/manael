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

// Package manael provides HTTP handler for processing images.
package manael_test // import "manael.org/x/manael"

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"manael.org/x/manael"
)

var manaelTests = []struct {
	host string
	path string
	statusCode int
}{
	{
		"Manael",
		"/logo.png",
		200,
	},
	{
		"manael.local",
		"/404.html",
		404,
	},
}

func TestNewServeProxy(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/logo.png")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	p := manael.NewServeProxy(ts.URL)

	for _, tc := range manaelTests {
		req := httptest.NewRequest(http.MethodGet, "https://"+tc.host+tc.path, nil)

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status code is %d, want %d", got, want)
		}
	}
}
