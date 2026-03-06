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
	"strconv"
	"strings"
)

// QualityOptions holds the encoding quality settings parsed from the q query
// parameter. A zero value for any field means "not specified"; the encoder
// falls back to its per-format default in that case.
type QualityOptions struct {
	// Default is a universal quality applied to any format not explicitly
	// overridden by a format-specific field.
	Default int
	// WebP is the encoding quality for WebP images (1–100).
	WebP int
	// AVIF is the encoding quality for AVIF images (1–100).
	AVIF int
}

// clampQuality clamps q to the valid range [1, 100].
func clampQuality(q int) int {
	if q < 1 {
		return 1
	}
	if q > 100 {
		return 100
	}
	return q
}

// forFormat returns the encoding quality to use for the given MIME type.
// The fallback value is used when neither a format-specific nor a universal
// quality has been configured. It is safe to call on a nil receiver.
func (q *QualityOptions) forFormat(mimeType string, fallback int) int {
	if q == nil {
		return fallback
	}
	switch mimeType {
	case "image/webp":
		if q.WebP > 0 {
			return q.WebP
		}
	case "image/avif":
		if q.AVIF > 0 {
			return q.AVIF
		}
	}
	if q.Default > 0 {
		return q.Default
	}
	return fallback
}

// parseQualityParams parses the q query parameter from vals and returns a
// *QualityOptions describing the requested encoding quality. It returns nil
// when the q parameter is absent or all tokens are invalid.
//
// Supported forms:
//
//	?q=80               universal quality for all formats
//	?q=webp:70          format-specific quality
//	?q=avif:60,webp:80  comma-separated format-specific values
//	?q=70,avif:50       combined universal and format-specific
//	?q=avif:60&q=webp:80 repeated q parameters
//
// Out-of-range values are clamped to [1, 100]. Unrecognised format names and
// non-integer values are silently ignored.
func parseQualityParams(vals url.Values) *QualityOptions {
	qs, ok := vals["q"]
	if !ok || len(qs) == 0 {
		return nil
	}

	opts := &QualityOptions{}
	found := false

	for _, v := range qs {
		for _, token := range strings.Split(v, ",") {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			if idx := strings.IndexByte(token, ':'); idx >= 0 {
				format := strings.ToLower(strings.TrimSpace(token[:idx]))
				qStr := strings.TrimSpace(token[idx+1:])
				n, err := strconv.Atoi(qStr)
				if err != nil {
					continue
				}
				switch format {
				case "webp":
					opts.WebP = clampQuality(n)
					found = true
				case "avif":
					opts.AVIF = clampQuality(n)
					found = true
				}
			} else {
				n, err := strconv.Atoi(token)
				if err != nil {
					continue
				}
				opts.Default = clampQuality(n)
				found = true
			}
		}
	}

	if !found {
		return nil
	}
	return opts
}
