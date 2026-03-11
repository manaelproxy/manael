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

package httputil

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetVaryHeader(t *testing.T) {
	tests := []struct {
		name   string
		vary   string
		wantIn string
	}{
		{
			name:   "empty Vary header",
			vary:   "",
			wantIn: "Accept",
		},
		{
			name:   "Vary already contains Accept",
			vary:   "Accept",
			wantIn: "Accept",
		},
		{
			name:   "Vary contains other value",
			vary:   "Accept-Encoding",
			wantIn: "Accept",
		},
		{
			name:   "Vary contains Accept and other value",
			vary:   "Accept, Accept-Encoding",
			wantIn: "Accept",
		},
		{
			name:   "Vary contains Accept-Language",
			vary:   "Accept-Encoding, Accept-Language",
			wantIn: "Accept",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{
				Header: make(http.Header),
			}
			if tt.vary != "" {
				res.Header.Set("Vary", tt.vary)
			}

			SetVaryHeader(res)

			got := res.Header.Get("Vary")
			if got == "" {
				t.Fatal("Vary header should not be empty after SetVaryHeader")
			}
			if !strings.Contains(got, tt.wantIn) {
				t.Errorf("SetVaryHeader() = %q, want to contain %q", got, tt.wantIn)
			}
		})
	}
}

func TestUpdateContentDispositionFilename(t *testing.T) {
	tests := []struct {
		name    string
		cd      string
		typ     string
		wantCD  string
	}{
		{
			name:   "empty Content-Disposition",
			cd:     "",
			typ:    "image/webp",
			wantCD: "",
		},
		{
			name:   "attachment with jpeg filename converted to webp",
			cd:     `attachment; filename="photo.jpg"`,
			typ:    "image/webp",
			wantCD: `attachment; filename=photo.webp`,
		},
		{
			name:   "attachment with jpeg filename converted to avif",
			cd:     `attachment; filename="photo.jpg"`,
			typ:    "image/avif",
			wantCD: `attachment; filename=photo.avif`,
		},
		{
			name:   "attachment with png filename converted to webp",
			cd:     `attachment; filename="image.png"`,
			typ:    "image/webp",
			wantCD: `attachment; filename=image.webp`,
		},
		{
			name:   "unknown type leaves header unchanged",
			cd:     `attachment; filename="photo.jpg"`,
			typ:    "image/jpeg",
			wantCD: `attachment; filename="photo.jpg"`,
		},
		{
			name:   "no filename param leaves header unchanged",
			cd:     "attachment",
			typ:    "image/webp",
			wantCD: "attachment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{
				Header: make(http.Header),
			}
			if tt.cd != "" {
				res.Header.Set("Content-Disposition", tt.cd)
			}

			UpdateContentDispositionFilename(res, tt.typ)

			got := res.Header.Get("Content-Disposition")
			if tt.wantCD == "" {
				if got != "" {
					t.Errorf("UpdateContentDispositionFilename() set Content-Disposition = %q, want empty", got)
				}
				return
			}
			if got != tt.wantCD {
				t.Errorf("UpdateContentDispositionFilename() = %q, want %q", got, tt.wantCD)
			}
		})
	}
}

func TestScanAcceptHeader(t *testing.T) {
	tests := []struct {
		name        string
		accept      string
		contentType string
		enableAVIF  bool
		want        string
	}{
		{
			name:        "webp accepted",
			accept:      "image/webp,image/*,*/*;q=0.8",
			contentType: "image/jpeg",
			want:        "image/webp",
		},
		{
			name:        "avif accepted with avif enabled",
			accept:      "image/avif,image/webp,*/*",
			contentType: "image/jpeg",
			enableAVIF:  true,
			want:        "image/avif",
		},
		{
			name:        "avif accepted but avif disabled",
			accept:      "image/avif,image/webp,*/*",
			contentType: "image/jpeg",
			enableAVIF:  false,
			want:        "image/webp",
		},
		{
			name:        "no compatible format",
			accept:      "text/html,*/*;q=0.8",
			contentType: "image/jpeg",
			want:        "*/*",
		},
		{
			name:        "empty accept",
			accept:      "",
			contentType: "image/jpeg",
			want:        "*/*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept", tt.accept)

			res := &http.Response{
				Request: req,
				Header:  make(http.Header),
			}
			res.Header.Set("Content-Type", tt.contentType)

			got := ScanAcceptHeader(res, tt.enableAVIF)
			if got != tt.want {
				t.Errorf("ScanAcceptHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSelectOutputType(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		statusCode   int
		contentType  string
		cacheControl string
		accept       string
		enableAVIF   bool
		want         string
	}{
		{
			name:        "jpeg with webp accepted",
			method:      http.MethodGet,
			statusCode:  http.StatusOK,
			contentType: "image/jpeg",
			accept:      "image/webp,*/*;q=0.8",
			want:        "image/webp",
		},
		{
			name:        "non-GET method",
			method:      http.MethodPost,
			statusCode:  http.StatusOK,
			contentType: "image/jpeg",
			accept:      "image/webp",
			want:        "*/*",
		},
		{
			name:        "non-OK status",
			method:      http.MethodGet,
			statusCode:  http.StatusNotFound,
			contentType: "image/jpeg",
			accept:      "image/webp",
			want:        "*/*",
		},
		{
			name:         "no-transform cache control",
			method:       http.MethodGet,
			statusCode:   http.StatusOK,
			contentType:  "image/jpeg",
			cacheControl: "no-transform",
			accept:       "image/webp",
			want:         "*/*",
		},
		{
			name:        "non-image content type",
			method:      http.MethodGet,
			statusCode:  http.StatusOK,
			contentType: "text/html",
			accept:      "*/*",
			want:        "*/*",
		},
		{
			name:       "304 not modified is allowed",
			method:     http.MethodGet,
			statusCode: http.StatusNotModified,
			contentType: "image/jpeg",
			accept:     "image/webp",
			want:       "image/webp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqHeader := make(http.Header)
			reqHeader.Set("Accept", tt.accept)

			req := &http.Request{
				Method: tt.method,
				Header: reqHeader,
			}

			res := &http.Response{
				Request:    req,
				StatusCode: tt.statusCode,
				Header:     make(http.Header),
			}
			if tt.contentType != "" {
				res.Header.Set("Content-Type", tt.contentType)
			}
			if tt.cacheControl != "" {
				res.Header.Set("Cache-Control", tt.cacheControl)
			}

			got := SelectOutputType(res, tt.enableAVIF)
			if got != tt.want {
				t.Errorf("SelectOutputType() = %q, want %q", got, tt.want)
			}
		})
	}
}


