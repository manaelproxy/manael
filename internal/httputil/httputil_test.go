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

			if got := res.Header.Get("Vary"); got == "" {
				t.Error("Vary header should not be empty after SetVaryHeader")
			}
		})
	}
}
