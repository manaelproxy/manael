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
package manael_test

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"manael.org/x/manael/v2"
	"manael.org/x/manael/v2/internal/testutil"
)

var basicTests = []struct {
	path       string
	statusCode int
}{
	{
		"/logo.png",
		200,
	},
	{
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

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	p := manael.NewServeProxy(u)

	for _, tc := range basicTests {
		req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status code is %d, want %d", got, want)
		}
	}
}

var varyTests = []struct {
	path string
	accept string
	vary string
}{
	{
		"/logo.png",
		"image/webp",
		"Accept",
	},
	{
		"/logo2.png",
		"image/webp",
		"Accept, Origin, Accept-Encoding",
	},
	{
		"/logo2.png",
		"image/png",
		"Origin, Accept-Encoding",
	},
}

func TestNewServeProxy_vary(t *testing.T) {
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

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	p := manael.NewServeProxy(u)

	for _, tc := range varyTests {
		req := httptest.NewRequest(http.MethodGet, ts.URL+tc.path, nil)
		req.Header.Set("Accept", tc.accept)

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

func TestNewServeProxy_badGateway(t *testing.T) {
	u, err := url.Parse("http://missing.test")
	if err != nil {
		t.Error(err)
	}

	p := manael.NewServeProxy(u)

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

var convertTests = []struct {
	accept      string
	path        string
	statusCode  int
	contentType string
	format      string
}{
	{
		"image/*,*/*;q=0.8",
		"/logo.png",
		http.StatusOK,
		"image/png",
		"png",
	},
	{
		"image/webp,image/*,*/*;q=0.8",
		"/logo.png",
		http.StatusOK,
		"image/webp",
		"webp",
	},
	{
		"image/*,*/*",
		"/photo.jpeg",
		http.StatusOK,
		"image/jpeg",
		"jpeg",
	},
	{
		"image/webp,image/*,*/*;q=0.8",
		"/photo.jpeg",
		http.StatusOK,
		"image/webp",
		"webp",
	},
	{
		"image/*,*/*;q=0.8",
		"/empty.gif",
		http.StatusOK,
		"image/gif",
		"gif",
	},
	{
		"image/webp,image/*,*/*;q=0.8",
		"/empty.gif",
		http.StatusOK,
		"image/gif",
		"gif",
	},
	{
		"image/webp,image/*,*/*",
		"/empty.txt",
		http.StatusOK,
		"text/plain; charset=utf-8",
		"not image",
	},
	{
		"image/webp,image/*,*/*",
		"/invalid.png",
		http.StatusOK,
		"image/png",
		"not image",
	},
}

func TestNewServeProxy_convert(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/logo.png")
	})
	mux.HandleFunc("/photo.jpeg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/photo.jpeg")
	})
	mux.HandleFunc("/empty.gif", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/empty.gif")
	})
	mux.HandleFunc("/empty.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/empty.txt")
	})
	mux.HandleFunc("/invalid.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/invalid.png")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	p := manael.NewServeProxy(u)

	for _, tc := range convertTests {
		req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)
		req.Header.Set("Accept", tc.accept)

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

		if got, want := testutil.DetectFormat(resp.Body), tc.format; got != want {
			t.Errorf("Detect format is %s, want %s", got, want)
		}
	}
}

var ifModifiedSinceTests = []struct {
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
		4090,
	},
	{
		"/logo.png",
		time.Date(2018, time.June, 30, 14, 3, 31, 0, time.UTC),
		http.StatusOK,
		4090,
	},
	{
		"/logo.png",
		time.Date(2018, time.June, 30, 14, 5, 31, 0, time.UTC),
		http.StatusNotModified,
		0,
	},
}

func TestNewServeProxy_ifModifiedSince(t *testing.T) {
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

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	p := manael.NewServeProxy(u)

	for _, tc := range ifModifiedSinceTests {
		req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)

		if !tc.modtime.IsZero() {
			req.Header.Set("If-Modified-Since", tc.modtime.Format(http.TimeFormat))
		}

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status Code is %d, want %d (%s)", got, want, tc.modtime)
		}

		body, _ := io.ReadAll(resp.Body)

		if got, want := len(body), tc.contentLength; got != want {
			t.Errorf("Response body is %d bytes, want %d", got, want)
		}
	}
}

var ifNoneMatchTests = []struct {
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
		4090,
	},
	{
		"/logo.png",
		"invalidETag",
		http.StatusOK,
		4090,
	},
}

func TestNewServeProxy_ifNoneMatch(t *testing.T) {
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

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	p := manael.NewServeProxy(u)

	for _, tc := range ifNoneMatchTests {
		req := httptest.NewRequest(http.MethodGet, "https://manael.local"+tc.path, nil)

		if tc.etag != "" {
			req.Header.Set("If-None-Match", tc.etag)
		}

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status Code is %d, want %d", got, want)
		}

		body, _ := io.ReadAll(resp.Body)

		if got, want := len(body), tc.contentLength; got != want {
			t.Errorf("Response body is %d bytes, want %d", got, want)
		}
	}
}

var acceptRangesTests = []struct {
	path         string
	accept       string
	acceptRanges string
}{
	{
		"/logo.png",
		"image/webp,image/*,*/*;q=0.8",
		"",
	},
	{
		"/logo.png",
		"image/*,*/*;q=0.8",
		"bytes",
	},
}

func TestNewServeProxy_acceptRanges(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Accept-Ranges", "bytes")
		http.ServeFile(w, r, "testdata/logo.png")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	p := manael.NewServeProxy(u)

	for _, tc := range acceptRangesTests {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		req.Header.Set("Accept", tc.accept)

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.Header.Get("Accept-Ranges"), tc.acceptRanges; got != want {
			t.Errorf(`Accept-Ranges is "%s", want "%s"`, got, want)
		}
	}
}

var noTransformTests = []struct {
	path        string
	contentType string
	format      string
}{
	{
		"/logo.png",
		"image/webp",
		"webp",
	},
	{
		"/logo.png?raw=1",
		"image/png",
		"png",
	},
	{
		"/photo.jpeg",
		"image/webp",
		"webp",
	},
	{
		"/photo.jpeg?raw=1",
		"image/jpeg",
		"jpeg",
	},
}

func TestNewServeProxy_noTransform(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("raw") == "1" {
			w.Header().Set("Cache-Control", "no-transform")
		}

		http.ServeFile(w, r, "testdata/logo.png")
	})
	mux.HandleFunc("/photo.jpeg", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("raw") == "1" {
			w.Header().Set("Cache-Control", "no-transform")
		}

		http.ServeFile(w, r, "testdata/photo.jpeg")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	p := manael.NewServeProxy(u)

	for _, tc := range noTransformTests {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		req.Header.Set("Accept", "image/webp,image/*,*/*;q=0.8")

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.Header.Get("Content-Type"), tc.contentType; got != want {
			t.Errorf("Content-Type is %s, want %s", got, want)
		}

		if got, want := testutil.DetectFormat(resp.Body), tc.format; got != want {
			t.Errorf("Detect format is %s, want %s", got, want)
		}
	}
}

var xForwardedForTests = []struct {
	path string
	xff  string
}{
	{
		"/xff.txt",
		"203.0.113.6, 192.0.2.44",
	},
}

func TestNewServeProxy_xForwardedFor(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/xff.txt", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, r.Header.Get("X-Forwarded-For"))
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	p := manael.NewServeProxy(u)

	for _, tc := range xForwardedForTests {
		req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)
		req.Header.Set("X-Forwarded-For", tc.xff)

		w := httptest.NewRecorder()

		p.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		xff, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		remoteIP, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := string(xff), fmt.Sprintf("%s, %s", tc.xff, remoteIP); got != want {
			t.Errorf(`X-Forwarded-For is "%s", want "%s"`, got, want)
		}
	}
}

var avifTests = []struct {
	accept      string
	path        string
	statusCode  int
	contentType string
}{
	{
		"image/avif,image/webp,image/*,*/*;q=0.8",
		"/photo.jpeg",
		http.StatusOK,
		"image/avif",
	},
	{
		"image/avif,image/webp,image/*,*/*;q=0.8",
		"/logo.png",
		http.StatusOK,
		"image/webp",
	},
}

func TestNewServeProxy_avif(t *testing.T) {
	if os.Getenv("MANAEL_ENABLE_AVIF") != "true" {
		t.Skip("Skipping test when avif disabled.")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/logo.png")
	})
	mux.HandleFunc("/photo.jpeg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/photo.jpeg")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	p := manael.NewServeProxy(u)

	for _, tc := range avifTests {
		req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)
		req.Header.Set("Accept", tc.accept)

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
	}
}
