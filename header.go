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

package manael

import (
	"mime"
	"net/http"
	"path"
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

func avifEnabled(res *http.Response, opts *ProxyOptions) bool {
	contentType := res.Header.Get("Content-Type")

	return opts.EnableAVIF && contentType != "image/png" && contentType != "image/gif"
}

func scanAcceptHeader(res *http.Response, opts *ProxyOptions) string {
	accepts := res.Request.Header.Get("Accept")

	for _, v := range strings.Split(accepts, ",") {
		t := strings.TrimSpace(v)

		if avifEnabled(res, opts) && strings.HasPrefix(t, "image/avif") {
			return "image/avif"
		} else if strings.HasPrefix(t, "image/webp") {
			return "image/webp"
		}
	}

	return "*/*"
}

func check(res *http.Response, opts *ProxyOptions) string {
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

	if t != "image/jpeg" && t != "image/png" && t != "image/gif" {
		return "*/*"
	}

	return scanAcceptHeader(res, opts)
}

func updateContentDispositionFilename(res *http.Response, typ string) {
	cd := res.Header.Get("Content-Disposition")
	if cd == "" {
		return
	}

	var ext string
	switch typ {
	case "image/webp":
		ext = ".webp"
	case "image/avif":
		ext = ".avif"
	default:
		return
	}

	disposition, params, err := mime.ParseMediaType(cd)
	if err != nil {
		return
	}

	if filename, ok := params["filename"]; ok {
		params["filename"] = strings.TrimSuffix(filename, path.Ext(filename)) + ext
		res.Header.Set("Content-Disposition", mime.FormatMediaType(disposition, params))
	}
}
