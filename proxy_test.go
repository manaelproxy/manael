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

package manael_test // import "manael.org/x/manael"

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"manael.org/x/manael"
)

var proxyTests = []struct {
	path string
	vary string
}{
	{
		"/logo.png",
		"Accept",
	},
	{
		"/logo2.png",
		"Accept, Origin, Accept-Encoding",
	},
}

func TestProxy_ServeHTTP(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "OriginServer")

		http.ServeFile(w, r, "testdata/logo.png")
	})
	mux.HandleFunc("/logo2.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Origin, Accept-Encoding")

		http.ServeFile(w, r, "testdata/logo.png")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	p := &manael.Proxy{http.DefaultTransport}

	for _, tc := range proxyTests {
		req := httptest.NewRequest(http.MethodGet, ts.URL+tc.path, nil)

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.Header.Get("Server"), "Manael"; got != want {
			t.Errorf("Server is %s, want %s", got, want)
		}

		if got, want := resp.Header.Get("Vary"), tc.vary; got != want {
			t.Errorf(`Vary is "%s", want "%s"`, got, want)
		}
	}
}

func TestProxy_ServeHTTP_badGateway(t *testing.T) {
	p := &manael.Proxy{http.DefaultTransport}

	req := httptest.NewRequest(http.MethodGet, "https://manael.invalid/test.png", nil)

	w := httptest.NewRecorder()

	p.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if got, want := resp.Header.Get("Server"), ""; got != want {
		t.Errorf(`Server is "%s", want "%s"`, got, want)
	}

	if got, want := resp.StatusCode, 502; got != want {
		t.Errorf("Status code is %d, want %d", got, want)
	}
}
