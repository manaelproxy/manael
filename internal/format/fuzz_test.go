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

package format

import (
	"bytes"
	"testing"
)

// FuzzIsAPNG verifies that IsAPNG does not panic for any byte input.
func FuzzIsAPNG(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	f.Add([]byte{'G', 'I', 'F', '8', '9', 'a'})

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = IsAPNG(bytes.NewReader(data))
	})
}

// FuzzIsAnimatedGIF verifies that IsAnimatedGIF does not panic for any byte input.
func FuzzIsAnimatedGIF(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte("GIF89a"))
	f.Add([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = IsAnimatedGIF(bytes.NewReader(data))
	})
}
