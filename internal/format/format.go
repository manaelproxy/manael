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

// Package format provides image format detection utilities.
package format

import (
	"bytes"
	"encoding/binary"
	"io"
)

var pngSignature = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

// IsAPNG returns true if r contains an APNG (Animated PNG) stream.
// It detects APNG by scanning PNG chunks for an acTL chunk before the IDAT chunk.
func IsAPNG(r io.Reader) (bool, error) {
	sig := make([]byte, 8)
	if _, err := io.ReadFull(r, sig); err != nil {
		return false, err
	}

	if !bytes.Equal(sig, pngSignature) {
		return false, nil
	}

	header := make([]byte, 8)
	for {
		if _, err := io.ReadFull(r, header); err != nil {
			return false, nil
		}

		chunkLen := binary.BigEndian.Uint32(header[0:4])
		chunkType := string(header[4:8])

		if chunkType == "acTL" {
			return true, nil
		}

		if chunkType == "IDAT" || chunkType == "IEND" {
			return false, nil
		}

		// Skip chunk data and CRC
		if _, err := io.CopyN(io.Discard, r, int64(chunkLen)+4); err != nil {
			return false, nil
		}
	}
}

// skipGIFSubBlocks discards GIF sub-block data until the block terminator (0x00).
func skipGIFSubBlocks(r io.Reader) error {
	b := make([]byte, 1)
	for {
		if _, err := io.ReadFull(r, b); err != nil {
			return err
		}
		if b[0] == 0 {
			return nil
		}
		if _, err := io.CopyN(io.Discard, r, int64(b[0])); err != nil {
			return err
		}
	}
}

// IsAnimatedGIF returns true if r contains an animated GIF stream (more than one image frame).
func IsAnimatedGIF(r io.Reader) (bool, error) {
	sig := make([]byte, 6)
	if _, err := io.ReadFull(r, sig); err != nil {
		return false, err
	}
	if string(sig[:3]) != "GIF" {
		return false, nil
	}

	// Read Logical Screen Descriptor (7 bytes).
	lsd := make([]byte, 7)
	if _, err := io.ReadFull(r, lsd); err != nil {
		return false, nil
	}

	// Skip Global Color Table if present.
	if lsd[4]&0x80 != 0 {
		size := 3 * (1 << (int(lsd[4]&0x07) + 1))
		if _, err := io.CopyN(io.Discard, r, int64(size)); err != nil {
			return false, nil
		}
	}

	frames := 0
	b := make([]byte, 1)
	for {
		if _, err := io.ReadFull(r, b); err != nil {
			break
		}
		switch b[0] {
		case 0x3B: // GIF Trailer
			return frames > 1, nil
		case 0x21: // Extension Introducer
			if _, err := io.ReadFull(r, b); err != nil {
				return false, nil
			}
			if err := skipGIFSubBlocks(r); err != nil {
				return false, nil
			}
		case 0x2C: // Image Descriptor
			frames++
			if frames > 1 {
				return true, nil
			}
			// Skip image descriptor fields (9 bytes: left, top, width, height, flags).
			desc := make([]byte, 9)
			if _, err := io.ReadFull(r, desc); err != nil {
				return false, nil
			}
			// Skip Local Color Table if present.
			if desc[8]&0x80 != 0 {
				size := 3 * (1 << (int(desc[8]&0x07) + 1))
				if _, err := io.CopyN(io.Discard, r, int64(size)); err != nil {
					return false, nil
				}
			}
			// Skip LZW minimum code size byte.
			if _, err := io.ReadFull(r, b); err != nil {
				return false, nil
			}
			// Skip image data sub-blocks.
			if err := skipGIFSubBlocks(r); err != nil {
				return false, nil
			}
		}
	}
	return frames > 1, nil
}
