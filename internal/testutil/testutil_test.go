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

// Package manael provides HTTP handler for processing images.
package testutil_test

import (
	"os"
	"testing"

	"manael.org/x/manael/v2/internal/testutil"
)

var testutilTests = []struct {
	name   string
	format string
}{
	{
		"../../testdata/logo.png",
		"png",
	},
	{
		"../../testdata/empty.txt",
		"not image",
	},
}

func TestDetectFormat(t *testing.T) {
	for _, tc := range testutilTests {
		f, err := os.Open(tc.name)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		format := testutil.DetectFormat(f)

		if got, want := format, tc.format; got != want {
			t.Errorf("A format is %s, want %s", got, want)
		}
	}
}
