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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var requestTests = []struct {
	path       string
	statusCode int
}{
	{
		"/logo.png",
		http.StatusOK,
	},
	{
		"/404.png",
		http.StatusNotFound,
	},
}

func TestRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, tc := range requestTests {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)

		resp, err := request(fmt.Sprintf("%s%s", ts.URL, tc.path), req)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status code is %d, want %d", got, want)
		}
	}
}

var requestTests2 = []struct {
	path       string
	modtime    time.Time
	statusCode int
}{
	{
		"/logo.png",
		time.Date(2018, time.June, 30, 14, 4, 31, 0, time.UTC),
		http.StatusNotModified,
	},
	{
		"/logo.png",
		time.Time{},
		http.StatusOK,
	},
	{
		"/logo.png",
		time.Date(2018, time.June, 30, 14, 3, 31, 0, time.UTC),
		http.StatusOK,
	},
	{
		"/logo.png",
		time.Date(2018, time.June, 30, 14, 5, 31, 0, time.UTC),
		http.StatusNotModified,
	},
}

func TestRequest_ifModifiedSince(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		modtime := time.Date(2018, time.June, 30, 14, 4, 31, 0, time.UTC)

		ims := r.Header.Get("If-Modified-Since")
		t, _ := time.Parse(http.TimeFormat, ims)

		if t.IsZero() || t.Before(modtime) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotModified)
		}
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, tc := range requestTests2 {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		if !tc.modtime.IsZero() {
			req.Header.Add("If-Modified-Since", tc.modtime.Format(http.TimeFormat))
		}

		resp, err := request(fmt.Sprintf("%s%s", ts.URL, tc.path), req)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status code is %d, want %d", got, want)
		}
	}
}

var requestTests3 = []struct {
	path       string
	etag       string
	statusCode int
}{
	{
		"/logo.png",
		`W/"fcaec3a55087c997f24ba2a70383ed9b7607fd85f0ae2e0dccb5ec094c75f009"`,
		http.StatusNotModified,
	},
	{
		"/logo.png",
		"",
		http.StatusOK,
	},
	{
		"/logo.png",
		"invalidETag",
		http.StatusOK,
	},
}

func TestRequest_ifNoneMatch(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/logo.png", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("If-None-Match") != fmt.Sprintf(`W/"%x"`, sha256.Sum256([]byte("etag"))) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotModified)
		}
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, tc := range requestTests3 {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)

		if tc.etag != "" {
			req.Header.Add("If-None-Match", tc.etag)
		}

		resp, err := request(fmt.Sprintf("%s%s", ts.URL, tc.path), req)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status code is %d, want %d", got, want)
		}
	}
}
