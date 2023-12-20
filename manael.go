// Copyright (c) 2017 Yamagishi Kazutoshi <ykzts@desire.sh>
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
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func setVaryHeader(res *http.Response) {
	keys := []string{"Accept"}
	for _, v := range strings.Split(res.Header.Get("Vary"), ",") {
		v = strings.TrimSpace(v)

		if v != "" && !strings.EqualFold(v, "Accept") {
			keys = append(keys, v)
		}
	}

	res.Header.Set("Vary", strings.Join(keys[:], ", "))
}

func avifEnabled(res *http.Response) bool {
	contentType := res.Header.Get("Content-Type")

	return os.Getenv("MANAEL_ENABLE_AVIF") == "true" && contentType != "image/png"
}

func scanAcceptHeader(res *http.Response) string {
	accepts := res.Request.Header.Get("Accept")

	for _, v := range strings.Split(accepts, ",") {
		t := strings.TrimSpace(v)

		if avifEnabled(res) && strings.HasPrefix(t, "image/avif") {
			return "image/avif"
		} else if strings.HasPrefix(t, "image/webp") {
			return "image/webp"
		}
	}

	return "*/*"
}

func check(res *http.Response) string {
	if res.Request.Method != http.MethodGet && res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotModified {
		return "*/*"
	}

	if s := res.Header.Get("Cache-Control"); s != "" {
		for _, v := range strings.Split(s, ",") {
			if strings.TrimSpace(v) == "no-transform" {
				return "*/*"
			}
		}
	}

	t := res.Header.Get("Content-Type")

	if t != "image/jpeg" && t != "image/png" {
		return "*/*"
	}

	return scanAcceptHeader(res)
}

func convert(src io.Reader, t string) (*bytes.Buffer, error) {
	img, err := Decode(src)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	err = Encode(buf, img, t)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func modifyResponse(res *http.Response) error {
	res.Header.Set("Server", "Manael")

	setVaryHeader(res)

	typ := check(res)
	if typ == "*/*" {
		return nil
	}

	defer res.Body.Close()

	p := bytes.NewBuffer(nil)
	b := io.TeeReader(res.Body, p)

	buf, err := convert(b, typ)
	if err != nil {
		body := io.MultiReader(p, res.Body)

		res.Body = io.NopCloser(body)
		log.Printf("error: %v\n", err)

		return nil
	}

	res.Body = io.NopCloser(buf)

	res.Header.Set("Content-Type", typ)
	res.Header.Set("Content-Length", strconv.Itoa(buf.Len()))

	if res.Header.Get("Accept-Ranges") != "" {
		res.Header.Del("Accept-Ranges")
	}

	return nil
}

// NewServeProxy returns a new Proxy given a upstream URL
func NewServeProxy(u *url.URL) http.Handler {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.Out.Header["X-Forwarded-For"] = r.In.Header["X-Forwarded-For"]

			r.SetXForwarded()
			r.SetURL(u)
		},
		ModifyResponse: modifyResponse,
	}
}
