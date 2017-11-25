package main

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

func decodeJPEG(src io.Reader) (img image.Image, err error) {
	img, err = jpeg.DecodeIntoRGBA(src, &jpeg.DecoderOptions{})
	if err != nil {
		return nil, err
	}

	return img, nil
}

func decodePNG(src io.Reader) (img image.Image, err error) {
	img, err = png.Decode(src)
	if err != nil {
		return nil, err
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

	return nil, errors.New("Not supported PNG format")
}

func decode(src io.Reader, contentType string) (img image.Image, err error) {
	switch contentType {
	case "image/jpeg":
		return decodeJPEG(src)
	case "image/png":
		return decodePNG(src)
	default:
		return nil, fmt.Errorf("Unknown content type: %s", contentType)
	}
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

func convert(src io.Reader, contentType string) (buf *bytes.Buffer, err error) {
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
