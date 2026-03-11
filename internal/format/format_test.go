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

func TestIsAPNG(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    bool
		wantErr bool
	}{
		{
			name:    "empty input",
			data:    []byte{},
			wantErr: true,
		},
		{
			name: "not a PNG",
			data: []byte{'G', 'I', 'F', '8', '9', 'a', 0x00, 0x00},
		},
		{
			name: "PNG signature only, no chunks",
			data: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
		},
		{
			name: "PNG with truncated chunk header",
			data: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00},
		},
		{
			name: "PNG with IDAT chunk (no acTL)",
			data: []byte{
				0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
				0x00, 0x00, 0x00, 0x00, 'I', 'D', 'A', 'T',
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsAPNG(bytes.NewReader(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("IsAPNG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsAPNG() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAnimatedGIF(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    bool
		wantErr bool
	}{
		{
			name:    "empty input",
			data:    []byte{},
			wantErr: true,
		},
		{
			name: "not a GIF",
			data: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
		},
		{
			name: "GIF header only",
			data: []byte("GIF89a"),
		},
		{
			name: "GIF with no frames (trailer only)",
			data: []byte{
				'G', 'I', 'F', '8', '9', 'a',
				0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x3B,
			},
		},
		{
			name: "GIF with one frame",
			data: []byte{
				'G', 'I', 'F', '8', '9', 'a',
				0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x2C,
				0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00,
				0x02,
				0x00,
				0x3B,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsAnimatedGIF(bytes.NewReader(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("IsAnimatedGIF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsAnimatedGIF() = %v, want %v", got, tt.want)
			}
		})
	}
}
