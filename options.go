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

// ProxyOptions holds configuration for a proxy created by NewServeProxy.
type ProxyOptions struct {
	// EnableAVIF enables AVIF encoding for JPEG images when the client supports it.
	EnableAVIF bool
	// EnableResize enables on-the-fly image resizing via the w, h, and fit
	// query parameters. Disabled by default; set to true to opt in.
	EnableResize bool
	// MaxImageSize is the maximum upstream image size in bytes that will be
	// converted. Images larger than this limit are passed through unchanged.
	// Defaults to 20 MiB when zero.
	MaxImageSize int64
	// MaxResizeWidth is the maximum allowed value for the w query parameter.
	// A value of 0 means no limit.
	MaxResizeWidth int
	// MaxResizeHeight is the maximum allowed value for the h query parameter.
	// A value of 0 means no limit.
	MaxResizeHeight int
	// AllowedWidths is an optional whitelist of permitted values for the w
	// query parameter. When non-empty, only listed widths are accepted.
	AllowedWidths []int
	// AllowedHeights is an optional whitelist of permitted values for the h
	// query parameter. When non-empty, only listed heights are accepted.
	AllowedHeights []int
	// DefaultQuality is the server-level fallback encoding quality (1–100)
	// used when no q query parameter is present in the request, or when the
	// request specifies format-specific quality without a universal default.
	// A value of 0 preserves the built-in per-format defaults (90 for WebP,
	// 60 for AVIF).
	DefaultQuality int
}

// ProxyOption is a functional option for configuring a proxy.
type ProxyOption func(*ProxyOptions)

// WithAVIFEnabled returns a ProxyOption that enables or disables AVIF encoding.
func WithAVIFEnabled(enabled bool) ProxyOption {
	return func(o *ProxyOptions) {
		o.EnableAVIF = enabled
	}
}

// WithResizeEnabled returns a ProxyOption that enables or disables on-the-fly
// image resizing via the w, h, and fit query parameters. Resizing is disabled
// by default; pass true to opt in. When disabled, resize query parameters are
// silently ignored and the original-size image is converted as usual.
func WithResizeEnabled(enabled bool) ProxyOption {
	return func(o *ProxyOptions) {
		o.EnableResize = enabled
	}
}

// WithMaxImageSize returns a ProxyOption that sets the maximum image size (in
// bytes) that will be converted. Images whose size exceeds this limit are
// passed through without conversion.
func WithMaxImageSize(size int64) ProxyOption {
	return func(o *ProxyOptions) {
		o.MaxImageSize = size
	}
}

// WithMaxResizeWidth returns a ProxyOption that sets the maximum allowed value
// for the w query parameter. Requests with a larger width will be rejected with
// 400 Bad Request. A value of 0 disables the limit.
func WithMaxResizeWidth(w int) ProxyOption {
	return func(o *ProxyOptions) {
		o.MaxResizeWidth = w
	}
}

// WithMaxResizeHeight returns a ProxyOption that sets the maximum allowed value
// for the h query parameter. Requests with a larger height will be rejected
// with 400 Bad Request. A value of 0 disables the limit.
func WithMaxResizeHeight(h int) ProxyOption {
	return func(o *ProxyOptions) {
		o.MaxResizeHeight = h
	}
}

// WithAllowedWidths returns a ProxyOption that restricts the w query parameter
// to the provided whitelist of pixel values. Requests specifying a width not
// in the list will be rejected with 400 Bad Request. An empty slice disables
// the whitelist.
func WithAllowedWidths(widths []int) ProxyOption {
	return func(o *ProxyOptions) {
		o.AllowedWidths = append([]int(nil), widths...)
	}
}

// WithAllowedHeights returns a ProxyOption that restricts the h query
// parameter to the provided whitelist of pixel values. Requests specifying a
// height not in the list will be rejected with 400 Bad Request. An empty
// slice disables the whitelist.
func WithAllowedHeights(heights []int) ProxyOption {
	return func(o *ProxyOptions) {
		o.AllowedHeights = append([]int(nil), heights...)
	}
}

// WithDefaultQuality returns a ProxyOption that sets the server-level fallback
// encoding quality (1–100). It is used when the request does not supply a q
// query parameter, or when the request specifies format-specific quality values
// without a universal default. A value of 0 (the default) preserves the
// built-in per-format defaults (90 for WebP, 60 for AVIF). Values outside
// [1, 100] are clamped to the nearest bound.
func WithDefaultQuality(q int) ProxyOption {
	return func(o *ProxyOptions) {
		o.DefaultQuality = clampQuality(q)
	}
}
