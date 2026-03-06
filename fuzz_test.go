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

package manael

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// FuzzScanAcceptHeader verifies that scanAcceptHeader does not panic for any
// value of the Accept request header.
func FuzzScanAcceptHeader(f *testing.F) {
	f.Add("image/webp,image/*,*/*;q=0.8")
	f.Add("image/avif,image/webp,*/*")
	f.Add("image/avif,image/webp,image/apng,image/*,*/*;q=0.8")
	f.Add("*/*")
	f.Add("")
	f.Add("image/webp")
	f.Add("image/avif")
	f.Add("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	f.Fuzz(func(t *testing.T, acceptHeader string) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept", acceptHeader)

		res := &http.Response{
			Request: req,
			Header:  make(http.Header),
		}
		res.Header.Set("Content-Type", "image/jpeg")

		_ = scanAcceptHeader(res, &ProxyOptions{})
	})
}

// FuzzCheck verifies that check does not panic for any combination of request
// method, response status, Content-Type, Cache-Control, and Accept headers.
// The http.Request is constructed directly (not via httptest.NewRequest) so
// that the fuzzer can supply arbitrary method strings without triggering the
// validation performed by http.NewRequest.
func FuzzCheck(f *testing.F) {
	f.Add(http.MethodGet, http.StatusOK, "image/jpeg", "", "image/webp,*/*;q=0.8")
	f.Add(http.MethodGet, http.StatusOK, "image/png", "", "image/webp,*/*;q=0.8")
	f.Add(http.MethodGet, http.StatusOK, "image/gif", "", "image/webp,*/*;q=0.8")
	f.Add(http.MethodGet, http.StatusOK, "image/jpeg", "no-transform", "image/webp")
	f.Add(http.MethodGet, http.StatusOK, "text/html", "", "*/*")
	f.Add(http.MethodPost, http.StatusOK, "image/jpeg", "", "image/webp")
	f.Add(http.MethodGet, http.StatusNotFound, "image/jpeg", "", "image/webp")
	f.Add(http.MethodGet, http.StatusOK, "", "", "")

	f.Fuzz(func(t *testing.T, method string, statusCode int, contentType, cacheControl, acceptHeader string) {
		reqHeader := make(http.Header)
		reqHeader.Set("Accept", acceptHeader)

		req := &http.Request{
			Method: method,
			Header: reqHeader,
		}

		res := &http.Response{
			Request:    req,
			StatusCode: statusCode,
			Header:     make(http.Header),
		}
		if contentType != "" {
			res.Header.Set("Content-Type", contentType)
		}
		if cacheControl != "" {
			res.Header.Set("Cache-Control", cacheControl)
		}

		_ = check(res, &ProxyOptions{})
	})
}

// FuzzIsAPNG verifies that isAPNG does not panic for any byte sequence.
func FuzzIsAPNG(f *testing.F) {
	// Valid PNG signature with no APNG marker.
	f.Add([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	// Empty input.
	f.Add([]byte{})
	// Not a PNG.
	f.Add([]byte("GIF89a"))
	// PNG signature followed by truncated chunk header.
	f.Add([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00})
	// PNG signature with IDAT chunk type.
	f.Add([]byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x00, 'I', 'D', 'A', 'T',
	})

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = isAPNG(bytes.NewReader(data))
	})
}

// FuzzIsAnimatedGIF verifies that isAnimatedGIF does not panic for any byte
// sequence.
func FuzzIsAnimatedGIF(f *testing.F) {
	// Valid GIF89a header only.
	f.Add([]byte("GIF89a"))
	// Empty input.
	f.Add([]byte{})
	// Not a GIF.
	f.Add([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	// GIF header with minimal logical screen descriptor (no global color table).
	f.Add([]byte{
		'G', 'I', 'F', '8', '9', 'a',
		0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
	})
	// GIF with trailer byte only.
	f.Add([]byte{
		'G', 'I', 'F', '8', '9', 'a',
		0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x3B,
	})
	// GIF with one image descriptor.
	f.Add([]byte{
		'G', 'I', 'F', '8', '9', 'a',
		0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x2C,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00,
		0x02,
		0x00,
		0x3B,
	})

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = isAnimatedGIF(bytes.NewReader(data))
	})
}

// FuzzSetVaryHeader verifies that setVaryHeader does not panic for any value
// of the Vary response header.
func FuzzSetVaryHeader(f *testing.F) {
	f.Add("Accept")
	f.Add("Accept, Accept-Encoding")
	f.Add("")
	f.Add("Accept-Encoding, Accept-Language")

	f.Fuzz(func(t *testing.T, vary string) {
		res := &http.Response{
			Header: make(http.Header),
		}
		if vary != "" {
			res.Header.Set("Vary", vary)
		}

		setVaryHeader(res)

		if got := res.Header.Get("Vary"); got == "" {
			t.Error("Vary header should not be empty after setVaryHeader")
		}
	})
}
