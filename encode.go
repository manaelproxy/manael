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

package manael

import (
	"bytes"
	"errors"
	"io"

	"github.com/h2non/bimg"
)

// Encode converts the image data in src to the format specified by Content-Type t,
// writing the result to w. An optional ResizeOptions value may be provided as the
// fourth argument to resize and/or crop the image in the same bimg pipeline pass
// before encoding. At most one ResizeOptions value is used; any additional values
// are ignored.
func Encode(w io.Writer, src []byte, t string, resize ...*ResizeOptions) error {
	var opts bimg.Options

	switch t {
	case "image/webp":
		opts = bimg.Options{
			Type:    bimg.WEBP,
			Quality: 90,
		}
	case "image/avif":
		opts = bimg.Options{
			Type:    bimg.AVIF,
			Quality: 60,
			Speed:   8,
		}
	default:
		return errors.New("unsupported image type")
	}

	var r *ResizeOptions
	if len(resize) > 0 {
		r = resize[0]
	}

	if r != nil && (r.Width > 0 || r.Height > 0) {
		opts.Width = r.Width
		opts.Height = r.Height
		switch r.Fit {
		case "cover":
			// Crop to fill the target dimensions exactly.
			opts.Crop = true
		case "contain":
			// Scale to fit within the target box; allow upscaling smaller images.
			opts.Enlarge = true
		case "scale-down", "":
			// Scale down to fit within the target box; never upscale.
			// This is the bimg default when neither Crop nor Enlarge is set.
		}
	}

	out, err := bimg.NewImage(src).Process(opts)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, bytes.NewReader(out))
	return err
}
