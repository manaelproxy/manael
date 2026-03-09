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

package manael_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/h2non/bimg"
	"manael.org/x/manael/v3"
)

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
