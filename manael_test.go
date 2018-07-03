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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewServeProxy(t *testing.T) {
	ts := httptest.NewServer(nil)
	defer ts.Close()

	p := NewServeProxy(ts.URL)

	if got, want := p.UpstreamURL.String(), ts.URL; got != want {
		t.Errorf("Upstream URL is %s, want %s", got, want)
	}
}

var serveProxyTests = []struct {
	accept      string
	path        string
	statusCode  int
	contentType string
	checksum    string
}{
	{
		"image/*,*/*;q=0.8",
		"/logo.png",
		http.StatusOK,
		"image/png",
		"87209ba2999bf451589e6333d0699dfcf6fa7c07eda5f4aa39051787dca62f2a",
	},
	{
		"image/webp,image/*,*/*;q=0.8",
		"/logo.png",
		http.StatusOK,
		"image/webp",
		"7af3f597f7965426d92f0333d27cf39377c9415107b2ffa716f72b7fe28ba2e9",
	},
	{
		"image/*,*/*",
		"/empty.jpeg",
		http.StatusOK,
		"image/jpeg",
		"1a8b498fcc782ef585778f6cd29640d78f0d2a25371786dcd62b61867d382b94",
	},
	{
		"image/webp,image/*,*/*;q=0.8",
		"/empty.jpeg",
		http.StatusOK,
		"image/webp",
		"69e3fa438596fdb60d6dee6aea10a92ebe2f19d50c4b069a8aaa97aa06b1a255",
	},
	{
		"image/*,*/*;q=0.8",
		"/empty.gif",
		http.StatusOK,
		"image/gif",
		"a065920df8cc4016d67c3a464be90099c9d28ffe7c9e6ee3a18f257efc58cbd7",
	},
	{
		"image/webp,image/*,*/*;q=0.8",
		"/empty.gif",
		http.StatusOK,
		"image/gif",
		"a065920df8cc4016d67c3a464be90099c9d28ffe7c9e6ee3a18f257efc58cbd7",
	},
	{
		"image/webp,image/*,*/*",
		"/empty.txt",
		http.StatusOK,
		"text/plain; charset=utf-8",
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	},
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

	for _, tc := range serveProxyTests {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
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

var serveProxyTests2 = []struct {
	path          string
	modtime       time.Time
	statusCode    int
	contentLength int
}{
	{
		"/logo.png",
		time.Date(2018, time.June, 30, 14, 4, 31, 0, time.UTC),
		http.StatusNotModified,
		0,
	},
	{
		"/logo.png",
		time.Time{},
		http.StatusOK,
		6435,
	},
	{
		"/logo.png",
		time.Date(2018, time.June, 30, 14, 3, 31, 0, time.UTC),
		http.StatusOK,
		6435,
	},
	{
		"/logo.png",
		time.Date(2018, time.June, 30, 14, 5, 31, 0, time.UTC),
		http.StatusNotModified,
		0,
	},
}

func TestServeProxy_ServeHTTP_ifModifiedSince(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		modtime := time.Date(2018, time.June, 30, 14, 4, 31, 0, time.UTC)

		ims := r.Header.Get("If-Modified-Since")
		t, _ := time.Parse(http.TimeFormat, ims)

		if t.IsZero() || t.Before(modtime) {
			r.Header.Del("If-Modified-Since")
			http.ServeFile(w, r, "testdata/logo.png")
		} else {
			w.WriteHeader(http.StatusNotModified)
		}
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	p := NewServeProxy(ts.URL)

	for _, tc := range serveProxyTests2 {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)

		if !tc.modtime.IsZero() {
			req.Header.Add("If-Modified-Since", tc.modtime.Format(http.TimeFormat))
		}

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status Code is %d, want %d (%s)", got, want, tc.modtime)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		if got, want := len(body), tc.contentLength; got != want {
			t.Errorf("Response body is %d bytes, want %d", got, want)
		}
	}
}

var serveProxyTests3 = []struct {
	path          string
	etag          string
	statusCode    int
	contentLength int
}{
	{
		"/logo.png",
		`W/"fcaec3a55087c997f24ba2a70383ed9b7607fd85f0ae2e0dccb5ec094c75f009"`,
		http.StatusNotModified,
		0,
	},
	{
		"/logo.png",
		"",
		http.StatusOK,
		6435,
	},
	{
		"/logo.png",
		"invalidETag",
		http.StatusOK,
		6435,
	},
}

func TestServeProxy_ServeHTTP_ifNoneMatch(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("If-None-Match") != fmt.Sprintf(`W/"%x"`, sha256.Sum256([]byte("etag"))) {
			http.ServeFile(w, r, "testdata/logo.png")
		} else {
			w.WriteHeader(http.StatusNotModified)
		}
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	p := NewServeProxy(ts.URL)

	for _, tc := range serveProxyTests3 {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)

		if tc.etag != "" {
			req.Header.Add("If-None-Match", tc.etag)
		}

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status Code is %d, want %d", got, want)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		if got, want := len(body), tc.contentLength; got != want {
			t.Errorf("Response body is %d bytes, want %d", got, want)
		}
	}
}
