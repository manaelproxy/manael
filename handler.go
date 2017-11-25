package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Handler is main process
type Handler struct {
	upstreamURL string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.RequestURI()
	url := h.upstreamURL + path

	resp, err := request(url, r)
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
