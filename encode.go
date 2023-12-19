// Copyright (c) 2021 Yamagishi Kazutoshi <ykzts@desire.sh>
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
	"errors"
	"image"
	"io"

	"github.com/Kagami/go-avif"
	"github.com/harukasan/go-libwebp/webp"
)

// Encode writes the Image m to w in the format specified by Content-Type t.
func Encode(w io.Writer, m image.Image, t string) error {
	switch t {
	case "image/avif":
		opts := avif.Options{
			Quality: 20,
			Speed:   8,
		}

		err := avif.Encode(w, m, &opts)
		if err != nil {
			return err
		}

		return nil
	case "image/webp":
		c, err := webp.ConfigPreset(webp.PresetDefault, 90)
		if err != nil {
			return err
		}

		switch img := m.(type) {
		case *image.Gray:
			err = webp.EncodeGray(w, img, c)
			if err != nil {
				return err
			}
		case *image.RGBA, *image.NRGBA:
			err = webp.EncodeRGBA(w, img, c)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return errors.New("Not supported image type")
}
