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
	"net/url"
	"testing"
)

var parseQualityParamsTests = []struct {
	name    string
	query   string
	wantNil bool
	want    QualityOptions
}{
	{
		name:    "no q parameter",
		query:   "",
		wantNil: true,
	},
	{
		name:  "universal quality",
		query: "q=80",
		want:  QualityOptions{Default: 80},
	},
	{
		name:  "webp-specific quality",
		query: "q=webp:70",
		want:  QualityOptions{WebP: 70},
	},
	{
		name:  "avif-specific quality",
		query: "q=avif:60",
		want:  QualityOptions{AVIF: 60},
	},
	{
		name:  "comma-separated format-specific values",
		query: "q=avif:60,webp:80",
		want:  QualityOptions{AVIF: 60, WebP: 80},
	},
	{
		name:  "combined universal and format-specific",
		query: "q=70,avif:50",
		want:  QualityOptions{Default: 70, AVIF: 50},
	},
	{
		name:  "multiple q params for different formats",
		query: "q=avif:60&q=webp:80",
		want:  QualityOptions{AVIF: 60, WebP: 80},
	},
	{
		name:  "multiple q params combined universal and format-specific",
		query: "q=70&q=avif:50",
		want:  QualityOptions{Default: 70, AVIF: 50},
	},
	{
		name:  "clamp below minimum to 1",
		query: "q=0",
		want:  QualityOptions{Default: 1},
	},
	{
		name:  "clamp above maximum to 100",
		query: "q=150",
		want:  QualityOptions{Default: 100},
	},
	{
		name:  "clamp negative value to 1",
		query: "q=-10",
		want:  QualityOptions{Default: 1},
	},
	{
		name:  "clamp format-specific value below minimum",
		query: "q=webp:0",
		want:  QualityOptions{WebP: 1},
	},
	{
		name:  "clamp format-specific value above maximum",
		query: "q=avif:200",
		want:  QualityOptions{AVIF: 100},
	},
	{
		name:    "invalid q value (non-integer)",
		query:   "q=abc",
		wantNil: true,
	},
	{
		name:    "invalid format-specific value (non-integer)",
		query:   "q=webp:abc",
		wantNil: true,
	},
	{
		name:    "unknown format name is ignored",
		query:   "q=jpeg:80",
		wantNil: true,
	},
	{
		name:    "empty q parameter value",
		query:   "q=",
		wantNil: true,
	},
}

func TestParseQualityParams(t *testing.T) {
	for _, tc := range parseQualityParamsTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			vals, err := url.ParseQuery(tc.query)
			if err != nil {
				t.Fatalf("url.ParseQuery(%q): %v", tc.query, err)
			}

			got := parseQualityParams(vals)

			if tc.wantNil {
				if got != nil {
					t.Errorf("parseQualityParams(%q) = %+v, want nil", tc.query, got)
				}
				return
			}

			if got == nil {
				t.Fatalf("parseQualityParams(%q) = nil, want %+v", tc.query, tc.want)
			}

			if got.Default != tc.want.Default {
				t.Errorf("parseQualityParams(%q).Default = %d, want %d", tc.query, got.Default, tc.want.Default)
			}
			if got.WebP != tc.want.WebP {
				t.Errorf("parseQualityParams(%q).WebP = %d, want %d", tc.query, got.WebP, tc.want.WebP)
			}
			if got.AVIF != tc.want.AVIF {
				t.Errorf("parseQualityParams(%q).AVIF = %d, want %d", tc.query, got.AVIF, tc.want.AVIF)
			}
		})
	}
}

var clampQualityTests = []struct {
	input int
	want  int
}{
	{0, 1},
	{1, 1},
	{50, 50},
	{100, 100},
	{101, 100},
	{-5, 1},
	{200, 100},
}

func TestClampQuality(t *testing.T) {
	for _, tc := range clampQualityTests {
		if got := clampQuality(tc.input); got != tc.want {
			t.Errorf("clampQuality(%d) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

var forFormatTests = []struct {
	name     string
	opts     *QualityOptions
	mimeType string
	fallback int
	want     int
}{
	{
		name:     "nil opts returns fallback",
		opts:     nil,
		mimeType: "image/webp",
		fallback: 90,
		want:     90,
	},
	{
		name:     "universal quality used for WebP when no format-specific",
		opts:     &QualityOptions{Default: 80},
		mimeType: "image/webp",
		fallback: 90,
		want:     80,
	},
	{
		name:     "WebP format-specific quality overrides universal",
		opts:     &QualityOptions{Default: 80, WebP: 70},
		mimeType: "image/webp",
		fallback: 90,
		want:     70,
	},
	{
		name:     "AVIF format-specific quality overrides universal",
		opts:     &QualityOptions{Default: 80, AVIF: 50},
		mimeType: "image/avif",
		fallback: 60,
		want:     50,
	},
	{
		name:     "universal quality used for AVIF when no format-specific",
		opts:     &QualityOptions{Default: 75},
		mimeType: "image/avif",
		fallback: 60,
		want:     75,
	},
	{
		name:     "fallback used when no quality specified at all",
		opts:     &QualityOptions{},
		mimeType: "image/webp",
		fallback: 90,
		want:     90,
	},
}

func TestQualityOptionsForFormat(t *testing.T) {
	for _, tc := range forFormatTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := tc.opts.forFormat(tc.mimeType, tc.fallback)
			if got != tc.want {
				t.Errorf("forFormat(%q, %d) = %d, want %d", tc.mimeType, tc.fallback, got, tc.want)
			}
		})
	}
}
