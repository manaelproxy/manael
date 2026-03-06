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
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// ResizeOptions holds the resize parameters parsed from request query parameters.
type ResizeOptions struct {
	// Width is the target width in pixels. Zero means auto (maintain aspect ratio).
	Width int
	// Height is the target height in pixels. Zero means auto (maintain aspect ratio).
	Height int
	// Fit controls how the image is resized to fit the target dimensions.
	// Supported values: "cover", "contain", "scale-down", or "" (default, behaves like contain).
	Fit string
}

// parseResizeParams parses the w, h, and fit query parameters from q and
// validates them against opts. It returns nil if no resize parameters are
// present. An error is returned if any parameter is invalid or violates
// configured limits.
func parseResizeParams(q url.Values, opts *ProxyOptions) (*ResizeOptions, error) {
	wStr := q.Get("w")
	hStr := q.Get("h")
	fit := q.Get("fit")

	if wStr == "" && hStr == "" && fit == "" {
		return nil, nil
	}

	var w, h int

	if wStr != "" {
		n, err := strconv.Atoi(wStr)
		if err != nil || n <= 0 {
			return nil, errors.New("invalid value for parameter 'w': must be a positive integer")
		}
		w = n
	}

	if hStr != "" {
		n, err := strconv.Atoi(hStr)
		if err != nil || n <= 0 {
			return nil, errors.New("invalid value for parameter 'h': must be a positive integer")
		}
		h = n
	}

	switch fit {
	case "", "cover", "contain", "scale-down":
	default:
		return nil, fmt.Errorf("invalid value for parameter 'fit': %q", fit)
	}

	if opts.MaxResizeWidth > 0 && w > opts.MaxResizeWidth {
		return nil, fmt.Errorf("width %d exceeds maximum allowed width %d", w, opts.MaxResizeWidth)
	}

	if opts.MaxResizeHeight > 0 && h > opts.MaxResizeHeight {
		return nil, fmt.Errorf("height %d exceeds maximum allowed height %d", h, opts.MaxResizeHeight)
	}

	if w > 0 && len(opts.AllowedWidths) > 0 {
		if !containsInt(opts.AllowedWidths, w) {
			return nil, fmt.Errorf("width %d is not in the list of allowed widths", w)
		}
	}

	if h > 0 && len(opts.AllowedHeights) > 0 {
		if !containsInt(opts.AllowedHeights, h) {
			return nil, fmt.Errorf("height %d is not in the list of allowed heights", h)
		}
	}

	return &ResizeOptions{Width: w, Height: h, Fit: fit}, nil
}

func containsInt(s []int, v int) bool {
	for _, n := range s {
		if n == v {
			return true
		}
	}
	return false
}
