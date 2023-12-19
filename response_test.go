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

package manael_test // import "manael.org/x/manael"

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"manael.org/x/manael"
)

var responseTests = []struct {
	statusCode    int
	status        string
	contentLength int64
	body          string
}{
	{
		http.StatusOK,
		"OK",
		6,
		"200 OK",
	},
	{
		http.StatusNotFound,
		"Not Found",
		13,
		"404 Not Found",
	},
	{
		http.StatusInternalServerError,
		"Internal Server Error",
		25,
		"500 Internal Server Error",
	},
}

func TestNewResponse(t *testing.T) {
	for _, tc := range responseTests {
		req := httptest.NewRequest(http.MethodGet, "https://manael.test/", nil)

		resp := manael.NewResponse(req, tc.statusCode)
		defer resp.Body.Close()

		if got, want := resp.StatusCode, tc.statusCode; got != want {
			t.Errorf("Status code is %d, want %d", got, want)
		}

		if got, want := resp.Status, tc.status; got != want {
			t.Errorf("Status is %s, want %s", got, want)
		}

		if got, want := resp.ContentLength, tc.contentLength; got != want {
			t.Errorf("Content length is %d, want %d", got, want)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := string(body), tc.body; got != want {
			t.Errorf("Response body is %s, want %s", got, want)
		}
	}
}
