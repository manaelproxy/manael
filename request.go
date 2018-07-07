// Copyright (c) 2018 Yamagishi Kazutoshi <ykzts@desire.sh>
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
	"net/http"
	"strings"
)

var client http.Client

func xff(r *http.Request) string {
	var ips []string

	if s := r.Header.Get("X-Forwarded-For"); s != "" {
		for _, ip := range strings.Split(s, ",") {
			ips = append(ips, strings.TrimSpace(ip))
		}
	}

	ips = append(ips, r.RemoteAddr)

	return strings.Join(ips[:], ",")
}

func request(url string, r *http.Request) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", r.Header.Get("User-Agent"))
	req.Header.Add("X-Forwarded-For", xff(r))

	for _, h := range []string{"X-Forwarded-Proto", "If-Modified-Since", "If-None-Match"} {
		v := r.Header.Get(h)
		if v != "" {
			req.Header.Add(h, v)
		}
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
