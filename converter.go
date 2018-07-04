// Copyright (c) 2017 Yamagishi Kazutoshi <ykzts@desire.sh>
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
	"errors"
	"image"
	// Register JPEG
	_ "image/jpeg"
	// Register PNG
	_ "image/png"
	"io"

	"github.com/harukasan/go-libwebp/webp"
)

func decode(src io.Reader) (image.Image, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	switch img.(type) {
	case *image.Gray, *image.RGBA, *image.NRGBA:
		return img, nil
	case *image.RGBA64, *image.NRGBA64, *image.YCbCr, *image.Paletted:
		bounds := img.Bounds()
		newImg := image.NewRGBA(bounds)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				newImg.Set(x, y, img.At(x, y))
			}
		}

		return newImg, nil
	}

	return nil, errors.New("Not supported image format")
}

func encode(src image.Image) (buf *bytes.Buffer, err error) {
	config, err := webp.ConfigPreset(webp.PresetDefault, 90)
	if err != nil {
		return nil, err
	}

	buf = bytes.NewBuffer(nil)
	switch img := src.(type) {
	case *image.Gray:
		err = webp.EncodeGray(buf, img, config)
		if err != nil {
			return nil, err
		}
	case *image.RGBA, *image.NRGBA:
		err = webp.EncodeRGBA(buf, img, config)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func convert(src io.Reader) (buf *bytes.Buffer, err error) {
	img, err := decode(src)
	if err != nil {
		return nil, err
	}

	buf, err = encode(img)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
