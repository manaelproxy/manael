// Copyright (c) 2019 Yamagishi Kazutoshi <ykzts@desire.sh>
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

package manael // import "manael.org/x/manael"

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

// Transport is an http.RoundTripper.
type Transport struct {
	UpstreamURL string

	Base http.RoundTripper
}

func (t *Transport) makeRequest(r *http.Request) (*http.Request, error) {
	u, err := url.Parse(t.UpstreamURL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, r.URL.Path)
	u.RawQuery = r.URL.RawQuery

	r2, err := http.NewRequest(r.Method, u.String(), r.Body)
	if err != nil {
		return nil, err
	}

	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if xff := r2.Header["X-Forwarded-For"]; len(xff) > 0 {
			ip = fmt.Sprintf("%s, %s", strings.Join(xff, ", "), ip)
		}

		r2.Header.Set("X-Forwarded-For", ip)
	}

	return r2, nil
}

func shouldEncodeToWebP(resp *http.Response) bool {
	if s := resp.Header.Get("Cache-Control"); s != "" {
		for _, v := range strings.Split(s, ",") {
			if strings.TrimSpace(v) == "no-transform" {
				return false
			}
		}
	}

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified) {
		return false
	}

	contentType := resp.Header.Get("Content-Type")
	return contentType == "image/jpeg" || contentType == "image/png"
}

func canDecodeWebP(r *http.Request) bool {
	accepts := r.Header.Get("Accept")

	for _, v := range strings.Split(accepts, ",") {
		t := strings.TrimSpace(v)
		if strings.HasPrefix(t, "image/webp") {
			return true
		}
	}

	return false
}

// RoundTrip responds to an converted image.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2, err := t.makeRequest(req)
	if err != nil {
		return nil, err
	}

	resp, err := t.Base.RoundTrip(req2)
	if err != nil {
		return nil, err
	}

	if !(shouldEncodeToWebP(resp) && canDecodeWebP(req2)) {
		return resp, nil
	}

	defer resp.Body.Close()

	buf, err := convert(resp.Body)
	if err != nil {
		resp = NewResponse(req2, http.StatusInternalServerError)
		log.Printf("error: %v\n", err)

		return resp, nil
	}

	resp.Body = ioutil.NopCloser(buf)

	resp.Header.Set("Content-Type", "image/webp")
	resp.Header.Set("Content-Length", strconv.Itoa(buf.Len()))

	if resp.Header.Get("Accept-Ranges") != "" {
		resp.Header.Del("Accept-Ranges")
	}

	return resp, nil
}
