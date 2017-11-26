package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

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

var client http.Client

// Handler is main process
type Handler struct {
	upstreamURL string
}

func (h *Handler) request(r *http.Request) (resp *http.Response, err error) {
	url := fmt.Sprintf("%s%s", h.upstreamURL, r.URL.RequestURI())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", r.Header.Get("User-Agent"))
	req.Header.Add("If-Modified-Since", r.Header.Get("If-Modified-Since"))

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) shouldEncodeToWebP(resp *http.Response) bool {
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified) {
		return false
	}

	contentType := resp.Header.Get("Content-Type")
	return contentType == "image/jpeg" || contentType == "image/png"
}

func (h *Handler) canDecodeWebP(r *http.Request) bool {
	accepts := r.Header.Get("Accept")

	for _, v := range strings.Split(accepts, ",") {
		t := strings.TrimSpace(v)
		if strings.HasPrefix(t, "image/webp") {
			return true
		}
	}

	return false
}

func (h *Handler) transfer(w http.ResponseWriter, resp *http.Response) {
	for key := range w.Header() {
		w.Header().Del(key)
	}
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.Header().Set("Vary", "Accept")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// ErrorInternalServerError return 500
func (h *Handler) ErrorInternalServerError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// ErrorBadGateway returns 502
func (h *Handler) ErrorBadGateway(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, "Bad Gateway", http.StatusBadGateway)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := h.request(r)
	if err != nil {
		h.ErrorBadGateway(w, r)
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if !(h.shouldEncodeToWebP(resp) && h.canDecodeWebP(r)) {
		h.transfer(w, resp)
		return
	}

	if resp.StatusCode == http.StatusNotModified {
		w.Header().Set("Cache-Control", resp.Header.Get("Cache-Control"))
		w.WriteHeader(http.StatusNotModified)
		return
	}

	buf, err := convert(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		h.ErrorInternalServerError(w, r)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "image/webp")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.Header().Set("Cache-Control", resp.Header.Get("Cache-Control"))
	w.Header().Set("Last-Modified", resp.Header.Get("Last-Modified"))
	w.Header().Set("Vary", "Accept")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, buf)
}
