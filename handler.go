package manael // import "manael.org/x/manael"

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var client http.Client

// Handler is main process
type Handler struct {
	UpstreamURL string
}

func (h *Handler) request(r *http.Request) (resp *http.Response, err error) {
	url := fmt.Sprintf("%s%s", h.UpstreamURL, r.URL.RequestURI())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", r.Header.Get("User-Agent"))

	ims := r.Header.Get("If-Modified-Since")
	if ims != "" {
		req.Header.Add("If-Modified-Since", ims)
	}

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

func (h *Handler) errorInternalServerError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (h *Handler) errorBadGateway(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, "Bad Gateway", http.StatusBadGateway)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := h.request(r)
	if err != nil {
		h.errorBadGateway(w, r)
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

	buf, err := Convert(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		h.errorInternalServerError(w, r)
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
