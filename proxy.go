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

// Package manael provides HTTP handler for processing images.
package manael

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"manael.org/x/manael/v3/internal/format"
	manaelhttputil "manael.org/x/manael/v3/internal/httputil"
)

// ctxKey is an unexported type for context keys used in this package.
type ctxKey int

const (
	resizeCtxKey  ctxKey = 0
	qualityCtxKey ctxKey = 1
)

func convert(src io.Reader, t string, resize *ResizeOptions, quality *QualityOptions) (*bytes.Buffer, error) {
	data, err := Decode(src)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	err = Encode(buf, data, t, resize, quality)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func modifyResponse(res *http.Response, opts *ProxyOptions) error {
	res.Header.Set("Server", "Manael")

	manaelhttputil.SetVaryHeader(res)

	typ := manaelhttputil.Check(res, opts.EnableAVIF)
	if typ == "*/*" {
		return nil
	}

	// Retrieve resize and quality options stored in the request context.
	resize, _ := res.Request.Context().Value(resizeCtxKey).(*ResizeOptions)
	quality, _ := res.Request.Context().Value(qualityCtxKey).(*QualityOptions)

	// Apply the server-level default quality when the request does not supply
	// a universal quality value (but may still have format-specific overrides).
	if opts.DefaultQuality > 0 {
		if quality == nil {
			quality = &QualityOptions{Default: opts.DefaultQuality}
		} else if quality.Default == 0 {
			merged := *quality
			merged.Default = opts.DefaultQuality
			quality = &merged
		}
	}

	maxSize := opts.MaxImageSize

	// If Content-Length is known and exceeds the limit, pass through unchanged.
	if res.ContentLength > maxSize {
		return nil
	}

	origBody := res.Body
	closeOrigBody := true
	defer func() {
		if closeOrigBody {
			_ = origBody.Close()
		}
	}()

	p := bytes.NewBuffer(nil)
	// Read at most maxSize+1 bytes so we can detect oversized streaming responses.
	limited := io.LimitReader(origBody, maxSize+1)
	b := io.TeeReader(limited, p)

	if res.Header.Get("Content-Type") == "image/png" {
		ok, _ := format.IsAPNG(b)
		// Drain remaining bytes into p so the full read is buffered.
		if _, err := io.Copy(io.Discard, b); err != nil {
			return err
		}
		if ok {
			// APNG: pass through unchanged (animated, no conversion).
			// Reconstruct with origBody remainder in case body exceeded maxSize.
			res.Body = struct {
				io.Reader
				io.Closer
			}{
				Reader: io.MultiReader(bytes.NewReader(p.Bytes()), origBody),
				Closer: origBody,
			}
			closeOrigBody = false
			return nil
		}
		// If oversized, pass through unchanged.
		if int64(p.Len()) > maxSize {
			res.Body = struct {
				io.Reader
				io.Closer
			}{
				Reader: io.MultiReader(bytes.NewReader(p.Bytes()), origBody),
				Closer: origBody,
			}
			closeOrigBody = false
			return nil
		}
		// Not APNG: replace b with a reader over p's buffered content for conversion.
		b = bytes.NewReader(p.Bytes())
	}

	if res.Header.Get("Content-Type") == "image/gif" {
		ok, _ := format.IsAnimatedGIF(b)
		if ok {
			// Animated GIF: pass through unchanged without buffering full payload.
			res.Body = struct {
				io.Reader
				io.Closer
			}{
				Reader: io.MultiReader(bytes.NewReader(p.Bytes()), origBody),
				Closer: origBody,
			}
			closeOrigBody = false
			return nil
		}
		// Not animated: buffer full body for conversion.
		if _, err := io.Copy(io.Discard, b); err != nil {
			return err
		}
		// If oversized, pass through unchanged.
		if int64(p.Len()) > maxSize {
			res.Body = struct {
				io.Reader
				io.Closer
			}{
				Reader: io.MultiReader(bytes.NewReader(p.Bytes()), origBody),
				Closer: origBody,
			}
			closeOrigBody = false
			return nil
		}
		// Replace b with a reader over p's buffered content for conversion.
		b = bytes.NewReader(p.Bytes())
	}

	buf, err := convert(b, typ, resize, quality)
	if err != nil {
		res.Body = struct {
			io.Reader
			io.Closer
		}{
			Reader: io.MultiReader(bytes.NewReader(p.Bytes()), origBody),
			Closer: origBody,
		}
		closeOrigBody = false
		slog.Error("failed to convert image",
			slog.String("error", err.Error()),
			slog.String("url", res.Request.URL.String()),
		)

		return nil
	}

	res.Body = io.NopCloser(buf)

	res.Header.Set("Content-Type", typ)
	res.Header.Set("Content-Length", strconv.Itoa(buf.Len()))

	manaelhttputil.UpdateContentDispositionFilename(res, typ)

	if res.Header.Get("Accept-Ranges") != "" {
		res.Header.Del("Accept-Ranges")
	}

	return nil
}

// NewServeProxy returns a new reverse-proxy handler for the given upstream URL.
// Optional ProxyOption values can be used to configure the proxy behaviour,
// such as enabling AVIF encoding, adjusting the maximum image size, or
// restricting image resize dimensions.
func NewServeProxy(u *url.URL, opts ...ProxyOption) http.Handler {
	options := &ProxyOptions{
		MaxImageSize: defaultMaxImageSize,
	}
	for _, o := range opts {
		o(options)
	}

	rp := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.Out.Header["X-Forwarded-For"] = r.In.Header["X-Forwarded-For"]

			r.SetXForwarded()
			r.SetURL(u)
		},
		ModifyResponse: func(res *http.Response) error {
			return modifyResponse(res, options)
		},
	}

	return &resizeProxy{proxy: rp, opts: options}
}

// resizeProxy wraps a ReverseProxy to validate and store resize parameters
// before forwarding the request upstream.
type resizeProxy struct {
	proxy *httputil.ReverseProxy
	opts  *ProxyOptions
}

// resizeQueryParams is the set of query parameter names consumed by Manael
// for resizing that must be stripped from the outbound upstream URL when
// resizing is enabled.
var resizeQueryParams = []string{"w", "h", "fit"}

// qualityQueryParams is the set of query parameter names consumed by Manael
// for quality control that must always be stripped from the outbound upstream URL.
var qualityQueryParams = []string{"q"}

func (s *resizeProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	if s.opts.EnableResize {
		resize, err := parseResizeParams(q, s.opts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if resize != nil {
			r = r.WithContext(context.WithValue(r.Context(), resizeCtxKey, resize))
		}

		// Strip resize params only when resizing is enabled and they are consumed.
		for _, p := range resizeQueryParams {
			q.Del(p)
		}
	}

	if quality := parseQualityParams(q); quality != nil {
		r = r.WithContext(context.WithValue(r.Context(), qualityCtxKey, quality))
	}

	// Quality params are always consumed by Manael; strip them before forwarding.
	for _, p := range qualityQueryParams {
		q.Del(p)
	}

	outURL := *r.URL
	outURL.RawQuery = q.Encode()
	r2 := r.Clone(r.Context())
	r2.URL = &outURL

	s.proxy.ServeHTTP(w, r2)
}
