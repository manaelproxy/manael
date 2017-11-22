package main

import (
	"flag"
	"image"
	"io"
	"log"
	"net/http"

	"github.com/harukasan/go-libwebp/webp"
	"github.com/pixiv/go-libjpeg/jpeg"
)

var port = flag.String("port", "8080", "")
var upstreamUrl = flag.String("upstream-url", "http://localhost:9000", "")

func supportsWebp(r *http.Request) bool {
	for _, value := range r.Header["Accept"] {
		if value == "image/webp" {
			return true
		}
	}
	return false
}

func decode(src io.Reader) (image.Image, error) {
	img, err := jpeg.DecodeIntoRGBA(src, &jpeg.DecoderOptions{})
	if err != nil {
		return nil, err
	}

	return img, nil
}

func encode(w io.Writer, img image.Image) (err error) {
	config, err := webp.ConfigPreset(webp.PresetDefault, 90)
	if err != nil {
		return err
	}

	if err = webp.EncodeRGBA(w, img, config); err != nil {
		return err
	}

	return nil
}

func transfer(w http.ResponseWriter, resp *http.Response) {
	for key, _ := range w.Header() {
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
	url := *upstreamUrl + path

	resp, err := http.Get(url)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header["Content-Type"][0]

	if contentType != "image/jpeg" || !supportsWebp(r) {
		transfer(w, resp)
		return
	}

	img, err := decode(resp.Body)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "image/webp")

	if err := encode(w, img); err != nil {
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
