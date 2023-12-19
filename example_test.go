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

package manael_test // import "manael.org/x/manael"

import (
	"fmt"
	"log"
	"net/http"

	"manael.org/x/manael"
)

func ExampleNewServeProxy() {
	p := manael.NewServeProxy("http://localhost:9000")

	err := http.ListenAndServe(":8080", p)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleNewResponse() {
	req, err := http.NewRequest(http.MethodGet, "https://manael.test/", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp := manael.NewResponse(req, http.StatusOK)

	fmt.Println(resp.StatusCode)
	// Output: 200
}
