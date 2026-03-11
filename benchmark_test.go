// Copyright (c) 2026 Yamagishi Kazutoshi <ykzts@desire.sh>
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
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"manael.org/x/manael/v3"
)

func BenchmarkEncode_PNGToWebP(b *testing.B) {
	src, err := os.ReadFile("testdata/logo.png")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		if err := manael.Encode(io.Discard, src, "image/webp", nil, nil); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncode_JPEGToWebP(b *testing.B) {
	src, err := os.ReadFile("testdata/photo.jpeg")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		if err := manael.Encode(io.Discard, src, "image/webp", nil, nil); err != nil {
			b.Fatal(err)
		}
	}
}

var benchmarkProxyTests = []struct {
	name   string
	file   string
	accept string
}{
	{
		name:   "PNGToWebP",
		file:   "testdata/logo.png",
		accept: "image/webp,image/*,*/*;q=0.8",
	},
	{
		name:   "JPEGToWebP",
		file:   "testdata/photo.jpeg",
		accept: "image/webp,image/*,*/*;q=0.8",
	},
	{
		name:   "PNGPassthrough",
		file:   "testdata/logo.png",
		accept: "image/*,*/*;q=0.8",
	},
}

func BenchmarkNewServeProxy(b *testing.B) {
	for _, tc := range benchmarkProxyTests {
		b.Run(tc.name, func(b *testing.B) {
			body, err := os.ReadFile(tc.file)
			if err != nil {
				b.Fatal(err)
			}

			mux := http.NewServeMux()
			mux.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", http.DetectContentType(body))
				_, _ = w.Write(body)
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			if err != nil {
				b.Fatal(err)
			}

			p := manael.NewServeProxy(u)

			b.ResetTimer()
			b.ReportAllocs()

			for b.Loop() {
				req := httptest.NewRequest(http.MethodGet, "https://manael.test/image", nil)
				req.Header.Set("Accept", tc.accept)

				w := httptest.NewRecorder()
				p.ServeHTTP(w, req)

				resp := w.Result()
				_, _ = io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	src, err := os.ReadFile("testdata/photo.jpeg")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		if _, err := manael.Decode(bytes.NewReader(src)); err != nil {
			b.Fatal(err)
		}
	}
}
