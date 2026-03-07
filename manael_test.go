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
	"strconv"
	"testing"
	"time"

	"github.com/h2non/bimg"
	"manael.org/x/manael/v3"
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
	path   string
	accept string
	vary   string
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
		"Accept, Origin, Accept-Encoding",
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
		"image/png",
	},
	{
		"image/webp,image/*,*/*;q=0.8",
		"/logo.png",
		http.StatusOK,
		"image/webp",
		"image/webp",
	},
	{
		"image/*,*/*",
		"/photo.jpeg",
		http.StatusOK,
		"image/jpeg",
		"image/jpeg",
	},
	{
		"image/webp,image/*,*/*;q=0.8",
		"/photo.jpeg",
		http.StatusOK,
		"image/webp",
		"image/webp",
	},
	{
		"image/*,*/*;q=0.8",
		"/empty.gif",
		http.StatusOK,
		"image/gif",
		"image/gif",
	},
	{
		"image/webp,image/*,*/*;q=0.8",
		"/empty.gif",
		http.StatusOK,
		"image/webp",
		"image/webp",
	},
	{
		"image/*,*/*;q=0.8",
		"/animation.gif",
		http.StatusOK,
		"image/gif",
		"image/gif",
	},
	{
		"image/webp,image/*,*/*;q=0.8",
		"/animation.gif",
		http.StatusOK,
		"image/gif",
		"image/gif",
	},
	{
		"image/*,*/*;q=0.8",
		"/animation.png",
		http.StatusOK,
		"image/png",
		"image/png",
	},
	{
		"image/webp,image/*,*/*;q=0.8",
		"/animation.png",
		http.StatusOK,
		"image/png",
		"image/png",
	},
	{
		"image/webp,image/*,*/*",
		"/empty.txt",
		http.StatusOK,
		"text/plain; charset=utf-8",
		"text/plain; charset=utf-8",
	},
	{
		"image/webp,image/*,*/*",
		"/invalid.png",
		http.StatusOK,
		"image/png",
		"text/plain; charset=utf-8",
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
	mux.HandleFunc("/animation.gif", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/animation.gif")
	})
	mux.HandleFunc("/animation.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/animation.png")
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		if got, want := http.DetectContentType(body), tc.format; got != want {
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
		"image/webp",
	},
	{
		"/logo.png?raw=1",
		"image/png",
		"image/png",
	},
	{
		"/photo.jpeg",
		"image/webp",
		"image/webp",
	},
	{
		"/photo.jpeg?raw=1",
		"image/jpeg",
		"image/jpeg",
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		if got, want := http.DetectContentType(body), tc.format; got != want {
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

var contentDispositionTests = []struct {
	name               string
	accept             string
	path               string
	file               string
	contentDisposition string
	want               string
}{
	{
		"WebP conversion with inline disposition",
		"image/webp,image/*,*/*;q=0.8",
		"/image",
		"testdata/logo.png",
		"inline; filename=logo.png",
		"inline; filename=logo.webp",
	},
	{
		"JPEG to WebP with quoted filename",
		"image/webp,image/*,*/*;q=0.8",
		"/image",
		"testdata/photo.jpeg",
		`attachment; filename="photo.jpeg"`,
		"attachment; filename=photo.webp",
	},
	{
		"No conversion when WebP not accepted",
		"image/*,*/*;q=0.8",
		"/image",
		"testdata/logo.png",
		"inline; filename=logo.png",
		"inline; filename=logo.png",
	},
	{
		"No Content-Disposition header",
		"image/webp,image/*,*/*;q=0.8",
		"/image",
		"testdata/logo.png",
		"",
		"",
	},
}

func TestNewServeProxy_contentDisposition(t *testing.T) {
	for _, tc := range contentDispositionTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(tc.path, func(w http.ResponseWriter, r *http.Request) {
				if tc.contentDisposition != "" {
					w.Header().Set("Content-Disposition", tc.contentDisposition)
				}
				http.ServeFile(w, r, tc.file)
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			if err != nil {
				t.Error(err)
			}

			p := manael.NewServeProxy(u)

			req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)
			req.Header.Set("Accept", tc.accept)

			w := httptest.NewRecorder()

			p.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if got, want := resp.Header.Get("Content-Disposition"), tc.want; got != want {
				t.Errorf(`Content-Disposition is %q, want %q`, got, want)
			}
		})
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
	t.Setenv("MANAEL_ENABLE_AVIF", "true")

	testImg, err := os.ReadFile("testdata/photo.jpeg")
	if err != nil {
		t.Fatal(err)
	}

	if err := manael.Encode(io.Discard, testImg, "image/avif", nil, nil); err != nil {
		t.Skipf("AVIF encoding not supported in this environment: %v", err)
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

	p := manael.NewServeProxy(u, manael.WithAVIFEnabled(true))

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

var maxImageSizeTests = []struct {
	name            string
	accept          string
	path            string
	file            string
	maxSize         string
	wantContentType string
	wantFormat      string
}{
	{
		// logo.png is 4090 bytes; limit of 1024 → pass through unchanged.
		"PNG exceeds limit, pass through",
		"image/webp,image/*,*/*;q=0.8",
		"/logo.png",
		"testdata/logo.png",
		"1024",
		"image/png",
		"image/png",
	},
	{
		// photo.jpeg is ~218 KB; limit of 1024 → pass through unchanged.
		"JPEG exceeds limit, pass through",
		"image/webp,image/*,*/*;q=0.8",
		"/photo.jpeg",
		"testdata/photo.jpeg",
		"1024",
		"image/jpeg",
		"image/jpeg",
	},
	{
		// logo.png is 4090 bytes; limit of 8192 → converted normally.
		"PNG within limit, convert",
		"image/webp,image/*,*/*;q=0.8",
		"/logo.png",
		"testdata/logo.png",
		"8192",
		"image/webp",
		"image/webp",
	},
}

// TestNewServeProxy_maxImageSize verifies that when WithMaxImageSize is set,
// images whose Content-Length exceeds the limit are passed through unconverted.
func TestNewServeProxy_maxImageSize(t *testing.T) {
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
		t.Fatal(err)
	}

	for _, tc := range maxImageSizeTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			maxSize, err := strconv.ParseInt(tc.maxSize, 10, 64)
			if err != nil {
				t.Fatalf("invalid maxSize %q: %v", tc.maxSize, err)
			}

			p := manael.NewServeProxy(u, manael.WithMaxImageSize(maxSize))

			req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)
			req.Header.Set("Accept", tc.accept)

			w := httptest.NewRecorder()
			p.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if got, want := resp.Header.Get("Content-Type"), tc.wantContentType; got != want {
				t.Errorf("Content-Type is %s, want %s", got, want)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := http.DetectContentType(body), tc.wantFormat; got != want {
				t.Errorf("Detected format is %s, want %s", got, want)
			}
		})
	}
}

// TestNewServeProxy_maxImageSizeStreaming verifies that oversized images served
// without a Content-Length header (chunked transfer) are also passed through
// unconverted when they exceed the configured MaxImageSize.
func TestNewServeProxy_maxImageSizeStreaming(t *testing.T) {
	// logo.png is 4090 bytes; limit of 1024 → pass through unchanged.
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("testdata/logo.png")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		// Flush headers before writing the body to force chunked transfer
		// encoding, which means no Content-Length header is sent.
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		w.Write(data) //nolint:errcheck
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	p := manael.NewServeProxy(u, manael.WithMaxImageSize(1024))

	req := httptest.NewRequest(http.MethodGet, "https://manael.test/logo.png", nil)
	req.Header.Set("Accept", "image/webp,image/*,*/*;q=0.8")

	w := httptest.NewRecorder()
	p.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if got, want := resp.Header.Get("Content-Type"), "image/png"; got != want {
		t.Errorf("Content-Type is %s, want %s", got, want)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := http.DetectContentType(body), "image/png"; got != want {
		t.Errorf("Detected format is %s, want %s", got, want)
	}
}

var resizeTests = []struct {
	name            string
	path            string
	accept          string
	wantStatusCode  int
	wantContentType string
	wantWidth       int
	wantHeight      int
}{
	{
		name:            "resize width only",
		path:            "/photo.jpeg?w=400",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
		wantWidth:       400,
		wantHeight:      300,
	},
	{
		name:            "resize height only",
		path:            "/photo.jpeg?h=300",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
		wantWidth:       400,
		wantHeight:      300,
	},
	{
		name:            "resize width and height with fit=cover",
		path:            "/photo.jpeg?w=200&h=200&fit=cover",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
		wantWidth:       200,
		wantHeight:      200,
	},
	{
		name:            "resize width and height with fit=contain",
		path:            "/photo.jpeg?w=200&h=200&fit=contain",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
		wantWidth:       200,
		wantHeight:      150,
	},
	{
		name:            "resize width and height with fit=scale-down",
		path:            "/photo.jpeg?w=200&h=200&fit=scale-down",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
		wantWidth:       200,
		wantHeight:      200,
	},
	{
		name:            "resize PNG image",
		path:            "/logo.png?w=128",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
		wantWidth:       128,
		wantHeight:      128,
	},
	{
		name:           "invalid w parameter (not a number)",
		path:           "/photo.jpeg?w=abc",
		accept:         "image/webp,image/*,*/*;q=0.8",
		wantStatusCode: http.StatusBadRequest,
	},
	{
		name:           "invalid w parameter (zero)",
		path:           "/photo.jpeg?w=0",
		accept:         "image/webp,image/*,*/*;q=0.8",
		wantStatusCode: http.StatusBadRequest,
	},
	{
		name:           "invalid h parameter (negative)",
		path:           "/photo.jpeg?h=-1",
		accept:         "image/webp,image/*,*/*;q=0.8",
		wantStatusCode: http.StatusBadRequest,
	},
	{
		name:           "invalid fit parameter",
		path:           "/photo.jpeg?w=200&fit=stretch",
		accept:         "image/webp,image/*,*/*;q=0.8",
		wantStatusCode: http.StatusBadRequest,
	},
}

// TestNewServeProxy_resize verifies that query-parameter based image resizing
// works correctly and that invalid parameters are rejected with 400 Bad Request.
func TestNewServeProxy_resize(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/photo.jpeg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/photo.jpeg")
	})
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/logo.png")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	p := manael.NewServeProxy(u, manael.WithResizeEnabled(true))

	for _, tc := range resizeTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)
			req.Header.Set("Accept", tc.accept)

			w := httptest.NewRecorder()
			p.ServeHTTP(w, req)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}

			if got, want := resp.StatusCode, tc.wantStatusCode; got != want {
				t.Errorf("Status Code is %d, want %d", got, want)
			}

			if tc.wantContentType != "" {
				if got, want := resp.Header.Get("Content-Type"), tc.wantContentType; got != want {
					t.Errorf("Content-Type is %s, want %s", got, want)
				}
			}

			if tc.wantWidth > 0 || tc.wantHeight > 0 {
				size, err := bimg.NewImage(body).Size()
				if err != nil {
					t.Fatalf("failed to get image size: %v", err)
				}
				if tc.wantWidth > 0 && size.Width != tc.wantWidth {
					t.Errorf("Width is %d, want %d", size.Width, tc.wantWidth)
				}
				if tc.wantHeight > 0 && size.Height != tc.wantHeight {
					t.Errorf("Height is %d, want %d", size.Height, tc.wantHeight)
				}
			}
		})
	}
}

var resizeLimitsTests = []struct {
	name           string
	path           string
	opts           []manael.ProxyOption
	wantStatusCode int
}{
	{
		name: "width within max limit",
		path: "/photo.jpeg?w=400",
		opts: []manael.ProxyOption{
			manael.WithMaxResizeWidth(4000),
		},
		wantStatusCode: http.StatusOK,
	},
	{
		name: "width exceeds max limit",
		path: "/photo.jpeg?w=5000",
		opts: []manael.ProxyOption{
			manael.WithMaxResizeWidth(4000),
		},
		wantStatusCode: http.StatusBadRequest,
	},
	{
		name: "height within max limit",
		path: "/photo.jpeg?h=300",
		opts: []manael.ProxyOption{
			manael.WithMaxResizeHeight(4000),
		},
		wantStatusCode: http.StatusOK,
	},
	{
		name: "height exceeds max limit",
		path: "/photo.jpeg?h=5000",
		opts: []manael.ProxyOption{
			manael.WithMaxResizeHeight(4000),
		},
		wantStatusCode: http.StatusBadRequest,
	},
	{
		name: "width in allowed widths whitelist",
		path: "/photo.jpeg?w=640",
		opts: []manael.ProxyOption{
			manael.WithAllowedWidths([]int{320, 640, 1280}),
		},
		wantStatusCode: http.StatusOK,
	},
	{
		name: "width not in allowed widths whitelist",
		path: "/photo.jpeg?w=500",
		opts: []manael.ProxyOption{
			manael.WithAllowedWidths([]int{320, 640, 1280}),
		},
		wantStatusCode: http.StatusBadRequest,
	},
	{
		name: "height in allowed heights whitelist",
		path: "/photo.jpeg?h=480",
		opts: []manael.ProxyOption{
			manael.WithAllowedHeights([]int{240, 480, 720}),
		},
		wantStatusCode: http.StatusOK,
	},
	{
		name: "height not in allowed heights whitelist",
		path: "/photo.jpeg?h=360",
		opts: []manael.ProxyOption{
			manael.WithAllowedHeights([]int{240, 480, 720}),
		},
		wantStatusCode: http.StatusBadRequest,
	},
	{
		name:           "no resize params, no limits applied",
		path:           "/photo.jpeg",
		opts:           []manael.ProxyOption{manael.WithMaxResizeWidth(100)},
		wantStatusCode: http.StatusOK,
	},
}

// TestNewServeProxy_resizeLimits verifies that MaxResizeWidth, MaxResizeHeight,
// AllowedWidths, and AllowedHeights configuration options are enforced correctly.
func TestNewServeProxy_resizeLimits(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/photo.jpeg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/photo.jpeg")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range resizeLimitsTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			opts := append([]manael.ProxyOption{manael.WithResizeEnabled(true)}, tc.opts...)
			p := manael.NewServeProxy(u, opts...)

			req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)
			req.Header.Set("Accept", "image/webp,image/*,*/*;q=0.8")

			w := httptest.NewRecorder()
			p.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if got, want := resp.StatusCode, tc.wantStatusCode; got != want {
				t.Errorf("Status Code is %d, want %d", got, want)
			}
		})
	}
}

// TestNewServeProxy_resizeDisabled verifies that when resize is not enabled
// (the default), resize query parameters are silently ignored and invalid
// resize parameters do not result in a 400 error.
func TestNewServeProxy_resizeDisabled(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/photo.jpeg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/photo.jpeg")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	// No WithResizeEnabled option — resize is disabled by default.
	p := manael.NewServeProxy(u)

	tests := []struct {
		name string
		path string
	}{
		{name: "w param ignored", path: "/photo.jpeg?w=100"},
		{name: "h param ignored", path: "/photo.jpeg?h=100"},
		{name: "w and h params ignored", path: "/photo.jpeg?w=200&h=200"},
		{name: "invalid w param ignored", path: "/photo.jpeg?w=abc"},
		{name: "negative h param ignored", path: "/photo.jpeg?h=-1"},
		{name: "invalid fit param ignored", path: "/photo.jpeg?w=200&fit=stretch"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)
			req.Header.Set("Accept", "image/webp,image/*,*/*;q=0.8")

			w := httptest.NewRecorder()
			p.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			// Resize params are ignored, so the request should succeed.
			if got, want := resp.StatusCode, http.StatusOK; got != want {
				t.Errorf("Status Code is %d, want %d (resize params should be ignored when disabled)", got, want)
			}
		})
	}
}

var qualityTests = []struct {
	name            string
	path            string
	accept          string
	opts            []manael.ProxyOption
	wantStatusCode  int
	wantContentType string
}{
	{
		name:            "universal quality via q param",
		path:            "/photo.jpeg?q=80",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
	},
	{
		name:            "webp-specific quality via q param",
		path:            "/photo.jpeg?q=webp:70",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
	},
	{
		name:            "comma-separated format-specific quality",
		path:            "/photo.jpeg?q=webp:70,avif:50",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
	},
	{
		name:            "combined universal and format-specific quality",
		path:            "/photo.jpeg?q=75,webp:65",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
	},
	{
		name:            "multiple q params",
		path:            "/photo.jpeg?q=webp:70&q=80",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
	},
	{
		name:            "out-of-range quality is clamped, not rejected",
		path:            "/photo.jpeg?q=200",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
	},
	{
		name:            "invalid q value falls back to default (no error)",
		path:            "/photo.jpeg?q=notanumber",
		accept:          "image/webp,image/*,*/*;q=0.8",
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
	},
	{
		name:  "WithDefaultQuality proxy option applied",
		path:  "/photo.jpeg",
		accept: "image/webp,image/*,*/*;q=0.8",
		opts: []manael.ProxyOption{
			manael.WithDefaultQuality(75),
		},
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
	},
	{
		name:  "q param overrides WithDefaultQuality",
		path:  "/photo.jpeg?q=65",
		accept: "image/webp,image/*,*/*;q=0.8",
		opts: []manael.ProxyOption{
			manael.WithDefaultQuality(75),
		},
		wantStatusCode:  http.StatusOK,
		wantContentType: "image/webp",
	},
}

// TestNewServeProxy_transformParamsNotForwarded verifies that Manael-specific
// transform query parameters (w, h, fit, q) are not forwarded to the upstream,
// while unrelated query parameters are preserved.
func TestNewServeProxy_transformParamsNotForwarded(t *testing.T) {
	var capturedRawQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/photo.jpeg", func(w http.ResponseWriter, r *http.Request) {
		capturedRawQuery = r.URL.RawQuery
		http.ServeFile(w, r, "testdata/photo.jpeg")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name              string
		path              string
		opts              []manael.ProxyOption
		wantMissingParams []string
		wantPresentParams map[string]string
	}{
		{
			name:              "resize params are stripped",
			path:              "/photo.jpeg?w=200&h=100&fit=cover",
			opts:              []manael.ProxyOption{manael.WithResizeEnabled(true)},
			wantMissingParams: []string{"w", "h", "fit"},
		},
		{
			name:              "quality param is stripped",
			path:              "/photo.jpeg?q=80",
			wantMissingParams: []string{"q"},
		},
		{
			name:              "all transform params stripped together",
			path:              "/photo.jpeg?w=200&h=100&fit=cover&q=75",
			opts:              []manael.ProxyOption{manael.WithResizeEnabled(true)},
			wantMissingParams: []string{"w", "h", "fit", "q"},
		},
		{
			name:              "non-transform params are preserved",
			path:              "/photo.jpeg?w=200&version=2&token=abc",
			opts:              []manael.ProxyOption{manael.WithResizeEnabled(true)},
			wantMissingParams: []string{"w"},
			wantPresentParams: map[string]string{"version": "2", "token": "abc"},
		},
		{
			name:              "only non-transform params passed unchanged",
			path:              "/photo.jpeg?foo=bar&baz=qux",
			wantMissingParams: []string{"w", "h", "fit", "q"},
			wantPresentParams: map[string]string{"foo": "bar", "baz": "qux"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			capturedRawQuery = ""

			p := manael.NewServeProxy(u, tc.opts...)

			req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)
			req.Header.Set("Accept", "image/webp,image/*,*/*;q=0.8")

			w := httptest.NewRecorder()
			p.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("unexpected status code %d", resp.StatusCode)
			}

			upstreamQ, err := url.ParseQuery(capturedRawQuery)
			if err != nil {
				t.Fatalf("failed to parse upstream query: %v", err)
			}

			for _, param := range tc.wantMissingParams {
				if _, present := upstreamQ[param]; present {
					t.Errorf("upstream received transform param %q, want it stripped", param)
				}
			}

			for param, want := range tc.wantPresentParams {
				if got := upstreamQ.Get(param); got != want {
					t.Errorf("upstream param %q = %q, want %q", param, got, want)
				}
			}
		})
	}
}

// TestNewServeProxy_quality verifies that the q query parameter controls
// encoding quality and that WithDefaultQuality is respected.
func TestNewServeProxy_quality(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/photo.jpeg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/photo.jpeg")
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range qualityTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := manael.NewServeProxy(u, tc.opts...)

			req := httptest.NewRequest(http.MethodGet, "https://manael.test"+tc.path, nil)
			req.Header.Set("Accept", tc.accept)

			w := httptest.NewRecorder()
			p.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if got, want := resp.StatusCode, tc.wantStatusCode; got != want {
				t.Errorf("Status Code is %d, want %d", got, want)
			}

			if tc.wantContentType != "" {
				if got, want := resp.Header.Get("Content-Type"), tc.wantContentType; got != want {
					t.Errorf("Content-Type is %s, want %s", got, want)
				}
			}
		})
	}
}
