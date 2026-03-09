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

package manael

import (
	"bytes"
	"testing"
)

// FuzzIsAPNG verifies that isAPNG does not panic for any byte sequence.
func FuzzIsAPNG(f *testing.F) {
	// Valid PNG signature with no APNG marker.
	f.Add([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	// Empty input.
	f.Add([]byte{})
	// Not a PNG.
	f.Add([]byte("GIF89a"))
	// PNG signature followed by truncated chunk header.
	f.Add([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00})
	// PNG signature with IDAT chunk type.
	f.Add([]byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x00, 'I', 'D', 'A', 'T',
	})

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = isAPNG(bytes.NewReader(data))
	})
}

// FuzzIsAnimatedGIF verifies that isAnimatedGIF does not panic for any byte
// sequence.
func FuzzIsAnimatedGIF(f *testing.F) {
	// Valid GIF89a header only.
	f.Add([]byte("GIF89a"))
	// Empty input.
	f.Add([]byte{})
	// Not a GIF.
	f.Add([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	// GIF header with minimal logical screen descriptor (no global color table).
	f.Add([]byte{
		'G', 'I', 'F', '8', '9', 'a',
		0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
	})
	// GIF with trailer byte only.
	f.Add([]byte{
		'G', 'I', 'F', '8', '9', 'a',
		0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x3B,
	})
	// GIF with one image descriptor.
	f.Add([]byte{
		'G', 'I', 'F', '8', '9', 'a',
		0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x2C,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00,
		0x02,
		0x00,
		0x3B,
	})

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = isAnimatedGIF(bytes.NewReader(data))
	})
}
