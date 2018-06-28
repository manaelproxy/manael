package manael // import "manael.org/x/manael"

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/harukasan/go-libwebp/webp"
	"github.com/pixiv/go-libjpeg/jpeg"
)

func decode(src io.Reader, contentType string) (img image.Image, err error) {
	switch contentType {
	case "image/jpeg":
		img, err = jpeg.DecodeIntoRGBA(src, &jpeg.DecoderOptions{})
	case "image/png":
		img, err = png.Decode(src)
	default:
		return nil, fmt.Errorf("Unknown content type: %s", contentType)
	}

	switch img.(type) {
	case *image.Gray, *image.RGBA, *image.NRGBA:
		return img, nil
	case *image.RGBA64, *image.NRGBA64, *image.Paletted:
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

func Convert(src io.Reader, contentType string) (buf *bytes.Buffer, err error) {
	img, err := decode(src, contentType)
	if err != nil {
		return nil, err
	}

	buf, err = encode(img)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
