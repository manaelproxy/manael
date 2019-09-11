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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewResponse returns response
func NewResponse(r *http.Request, status int) *http.Response {
	header := make(http.Header)
	statusText := http.StatusText(status)

	body := fmt.Sprintf("%d %s", status, statusText)
	buf := bytes.NewBufferString(body)

	header.Set("Content-Type", "text/plain; charset=utf-8")

	return &http.Response{
		Status:           statusText,
		StatusCode:       status,
		Header:           header,
		Body:             ioutil.NopCloser(buf),
		ContentLength:    int64(buf.Len()),
		Request:          r,
		TransferEncoding: r.TransferEncoding,
	}
}
