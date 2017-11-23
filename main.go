package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/harukasan/go-libwebp/webp"
	"github.com/pixiv/go-libjpeg/jpeg"
)

var port = flag.String("port", "8080", "")
var upstreamURL = flag.String("upstream-url", "http://localhost:9000", "")

func getMediaType(contentTypes []string) (format string) {
	contentType := contentTypes[len(contentTypes)-1]
	mediaType := strings.Split(contentType, ";")[0]
	return mediaType
}

func supportsWebp(accepts []string) bool {
	for _, v := range accepts {
		for _, contentType := range strings.Split(v, ",") {
			if strings.HasPrefix(contentType, "image/webp") {
				return true
			}
		}
	}
	return false
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

	return img, nil
}

func decode(src io.Reader, mediaType string) (img image.Image, err error) {
	switch mediaType {
	case "image/jpeg":
		return decodeJPEG(src)
	case "image/png":
		return decodePNG(src)
	default:
		return nil, fmt.Errorf("Unknown media type: %s", mediaType)
	}
}

func encode(w io.Writer, img image.Image) (err error) {
	config, err := webp.ConfigPreset(webp.PresetDefault, 90)
	if err != nil {
		return err
	}

	err = webp.EncodeRGBA(w, img, config)
	if err != nil {
		return err
	}

	return nil
}

func convert(w io.Writer, src io.Reader, mediaType string) (err error) {
	img, err := decode(src, mediaType)
	if err != nil {
		return err
	}

	err = encode(w, img)
	if err != nil {
		return err
	}

	return nil
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
	io.Copy(w, resp.Body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.RequestURI()
	url := *upstreamURL + path

	resp, err := http.Get(url)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	mediaType := getMediaType(resp.Header["Content-Type"])

	if !(mediaType == "image/jpeg" || mediaType == "image/png") || !supportsWebp(r.Header["Accept"]) {
		transfer(w, resp)
		return
	}

	w.Header().Set("Content-Type", "image/webp")

	err = convert(w, resp.Body, mediaType)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
