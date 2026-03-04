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

// Package manael provides HTTP handler for processing images.
package manael

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

var pngSignature = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

// isAPNG returns true if r contains an APNG (Animated PNG) stream.
// It detects APNG by scanning PNG chunks for an acTL chunk before the IDAT chunk.
func isAPNG(r io.Reader) (bool, error) {
	sig := make([]byte, 8)
	if _, err := io.ReadFull(r, sig); err != nil {
		return false, err
	}

	if !bytes.Equal(sig, pngSignature) {
		return false, nil
	}

	header := make([]byte, 8)
	for {
		if _, err := io.ReadFull(r, header); err != nil {
			return false, nil
		}

		chunkLen := binary.BigEndian.Uint32(header[0:4])
		chunkType := string(header[4:8])

		if chunkType == "acTL" {
			return true, nil
		}

		if chunkType == "IDAT" || chunkType == "IEND" {
			return false, nil
		}

		// Skip chunk data and CRC
		if _, err := io.CopyN(io.Discard, r, int64(chunkLen)+4); err != nil {
			return false, nil
		}
	}
}

// skipGIFSubBlocks discards GIF sub-block data until the block terminator (0x00).
func skipGIFSubBlocks(r io.Reader) error {
	b := make([]byte, 1)
	for {
		if _, err := io.ReadFull(r, b); err != nil {
			return err
		}
		if b[0] == 0 {
			return nil
		}
		if _, err := io.CopyN(io.Discard, r, int64(b[0])); err != nil {
			return err
		}
	}
}

// isAnimatedGIF returns true if r contains an animated GIF stream (more than one image frame).
func isAnimatedGIF(r io.Reader) (bool, error) {
	sig := make([]byte, 6)
	if _, err := io.ReadFull(r, sig); err != nil {
		return false, err
	}
	if string(sig[:3]) != "GIF" {
		return false, nil
	}

	// Read Logical Screen Descriptor (7 bytes).
	lsd := make([]byte, 7)
	if _, err := io.ReadFull(r, lsd); err != nil {
		return false, nil
	}

	// Skip Global Color Table if present.
	if lsd[4]&0x80 != 0 {
		size := 3 * (1 << (int(lsd[4]&0x07) + 1))
		if _, err := io.CopyN(io.Discard, r, int64(size)); err != nil {
			return false, nil
		}
	}

	frames := 0
	b := make([]byte, 1)
	for {
		if _, err := io.ReadFull(r, b); err != nil {
			break
		}
		switch b[0] {
		case 0x3B: // GIF Trailer
			return frames > 1, nil
		case 0x21: // Extension Introducer
			if _, err := io.ReadFull(r, b); err != nil {
				return false, nil
			}
			if err := skipGIFSubBlocks(r); err != nil {
				return false, nil
			}
		case 0x2C: // Image Descriptor
			frames++
			if frames > 1 {
				return true, nil
			}
			// Skip image descriptor fields (9 bytes: left, top, width, height, flags).
			desc := make([]byte, 9)
			if _, err := io.ReadFull(r, desc); err != nil {
				return false, nil
			}
			// Skip Local Color Table if present.
			if desc[8]&0x80 != 0 {
				size := 3 * (1 << (int(desc[8]&0x07) + 1))
				if _, err := io.CopyN(io.Discard, r, int64(size)); err != nil {
					return false, nil
				}
			}
			// Skip LZW minimum code size byte.
			if _, err := io.ReadFull(r, b); err != nil {
				return false, nil
			}
			// Skip image data sub-blocks.
			if err := skipGIFSubBlocks(r); err != nil {
				return false, nil
			}
		}
	}
	return frames > 1, nil
}

func setVaryHeader(res *http.Response) {
	keys := []string{"Accept"}
	for _, v := range strings.Split(res.Header.Get("Vary"), ",") {
		v = strings.TrimSpace(v)

		if v != "" && !strings.EqualFold(v, "Accept") {
			keys = append(keys, v)
		}
	}

	res.Header.Set("Vary", strings.Join(keys[:], ", "))
}

func avifEnabled(res *http.Response) bool {
	contentType := res.Header.Get("Content-Type")

	return os.Getenv("MANAEL_ENABLE_AVIF") == "true" && contentType != "image/png" && contentType != "image/gif"
}

func scanAcceptHeader(res *http.Response) string {
	accepts := res.Request.Header.Get("Accept")

	for _, v := range strings.Split(accepts, ",") {
		t := strings.TrimSpace(v)

		if avifEnabled(res) && strings.HasPrefix(t, "image/avif") {
			return "image/avif"
		} else if strings.HasPrefix(t, "image/webp") {
			return "image/webp"
		}
	}

	return "*/*"
}

func check(res *http.Response) string {
	if res.Request.Method != http.MethodGet && res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotModified {
		return "*/*"
	}

	if s := res.Header.Get("Cache-Control"); s != "" {
		for _, v := range strings.Split(s, ",") {
			if strings.TrimSpace(v) == "no-transform" {
				return "*/*"
			}
		}
	}

	t := res.Header.Get("Content-Type")

	if t != "image/jpeg" && t != "image/png" && t != "image/gif" {
		return "*/*"
	}

	return scanAcceptHeader(res)
}

func convert(src io.Reader, t string) (*bytes.Buffer, error) {
	img, err := Decode(src)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	err = Encode(buf, img, t)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func updateContentDispositionFilename(res *http.Response, typ string) {
	cd := res.Header.Get("Content-Disposition")
	if cd == "" {
		return
	}

	var ext string
	switch typ {
	case "image/webp":
		ext = ".webp"
	case "image/avif":
		ext = ".avif"
	default:
		return
	}

	disposition, params, err := mime.ParseMediaType(cd)
	if err != nil {
		return
	}

	if filename, ok := params["filename"]; ok {
		params["filename"] = strings.TrimSuffix(filename, path.Ext(filename)) + ext
		res.Header.Set("Content-Disposition", mime.FormatMediaType(disposition, params))
	}
}

func modifyResponse(res *http.Response) error {
	res.Header.Set("Server", "Manael")

	setVaryHeader(res)

	typ := check(res)
	if typ == "*/*" {
		return nil
	}

	origBody := res.Body
	closeOrigBody := true
	defer func() {
		if closeOrigBody {
			_ = origBody.Close()
		}
	}()

	p := bytes.NewBuffer(nil)
	b := io.TeeReader(origBody, p)

	if res.Header.Get("Content-Type") == "image/png" {
		ok, _ := isAPNG(b)
		// Drain remaining bytes into p so the full body is buffered.
		if _, err := io.Copy(io.Discard, b); err != nil {
			return err
		}
		if ok {
			// APNG: pass through unchanged; p now holds the complete body.
			res.Body = io.NopCloser(p)
			return nil
		}
		// Not APNG: replace b with a reader over p's buffered content for conversion.
		b = bytes.NewReader(p.Bytes())
	}

	if res.Header.Get("Content-Type") == "image/gif" {
		ok, _ := isAnimatedGIF(b)
		if ok {
			// Animated GIF: pass through unchanged without buffering full payload.
			res.Body = struct {
				io.Reader
				io.Closer
			}{
				Reader: io.MultiReader(bytes.NewReader(p.Bytes()), origBody),
				Closer: origBody,
			}
			closeOrigBody = false
			return nil
		}
		// Not animated: buffer full body for conversion.
		if _, err := io.Copy(io.Discard, b); err != nil {
			return err
		}
		// Replace b with a reader over p's buffered content for conversion.
		b = bytes.NewReader(p.Bytes())
	}

	buf, err := convert(b, typ)
	if err != nil {
		body := io.MultiReader(p, res.Body)

		res.Body = io.NopCloser(body)
		log.Printf("error: %v\n", err)

		return nil
	}

	res.Body = io.NopCloser(buf)

	res.Header.Set("Content-Type", typ)
	res.Header.Set("Content-Length", strconv.Itoa(buf.Len()))

	updateContentDispositionFilename(res, typ)

	if res.Header.Get("Accept-Ranges") != "" {
		res.Header.Del("Accept-Ranges")
	}

	return nil
}

// NewServeProxy returns a new Proxy given a upstream URL
func NewServeProxy(u *url.URL) http.Handler {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.Out.Header["X-Forwarded-For"] = r.In.Header["X-Forwarded-For"]

			r.SetXForwarded()
			r.SetURL(u)
		},
		ModifyResponse: modifyResponse,
	}
}
