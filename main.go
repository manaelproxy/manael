package main

import (
	"bytes"
	"errors"
	"flag"
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

var port = flag.String("port", "8080", "")
var upstreamURL = flag.String("upstream-url", "http://localhost:9000", "")

func request(url string, r *http.Request) (resp *http.Response, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("If-Modified-Since", r.Header.Get("If-Modified-Since"))

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func canDecodeWebP(r *http.Request) bool {
	accepts := r.Header.Get("Accept")

	for _, v := range strings.Split(accepts, ",") {
		t := strings.TrimSpace(v)
		if strings.HasPrefix(t, "image/webp") {
			return true
		}
	}

	return false
}

func shouldEncodeToWebP(resp *http.Response) bool {
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified) {
		return false
	}

	contentType := resp.Header.Get("Content-Type")
	return contentType == "image/jpeg" || contentType == "image/png"
}

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
	case *image.Paletted:
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

	buf = new(bytes.Buffer)
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

func transfer(w http.ResponseWriter, resp *http.Response) {
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

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.RequestURI()
	url := *upstreamURL + path

	resp, err := request(url, r)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if !(shouldEncodeToWebP(resp) && canDecodeWebP(r)) {
		transfer(w, resp)
		return
	}

	if resp.StatusCode == http.StatusNotModified {
		w.Header().Set("Cache-Control", resp.Header.Get("Cache-Control"))
		w.WriteHeader(http.StatusNotModified)
		return
	}

	buf, err := convert(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
