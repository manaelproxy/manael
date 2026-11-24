package manael_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/h2non/bimg"
	"manael.org/x/manael/v3"
)

// withOrientation returns jpeg with an EXIF APP1 segment that carries only
// the given orientation, inserted right after the SOI marker.
func withOrientation(jpeg []byte, orientation byte) []byte {
	tiff := []byte{
		'M', 'M', 0x00, 0x2a, 0x00, 0x00, 0x00, 0x08, // big-endian header, IFD at 8
		0x00, 0x01, // one entry
		0x01, 0x12, 0x00, 0x03, 0x00, 0x00, 0x00, 0x01, 0x00, orientation, 0x00, 0x00, // Orientation, SHORT
		0x00, 0x00, 0x00, 0x00, // no next IFD
	}
	payload := append([]byte("Exif\x00\x00"), tiff...)
	length := len(payload) + 2
	app1 := append([]byte{0xff, 0xe1, byte(length >> 8), byte(length)}, payload...)
	out := append([]byte{}, jpeg[:2]...)
	out = append(out, app1...)
	return append(out, jpeg[2:]...)
}

// A scale-down compares the box with the image as it is displayed: an 800×600
// JPEG rotated by its EXIF orientation is 600 wide, so a 700px box already
// holds it.
func TestEncode_scaleDownFollowsTheEXIFOrientation(t *testing.T) {
	src, err := os.ReadFile("testdata/photo.jpeg")
	if err != nil {
		t.Fatal(err)
	}
	rotated := withOrientation(src, 6)
	if meta, err := bimg.NewImage(rotated).Metadata(); err != nil || meta.Orientation != 6 {
		t.Fatalf("test image orientation = %d (%v), want 6", meta.Orientation, err)
	}

	var buf bytes.Buffer
	if err := manael.Encode(&buf, rotated, "image/webp", &manael.ResizeOptions{Width: 700, Fit: "scale-down"}, nil); err != nil {
		t.Fatal(err)
	}
	size, err := bimg.NewImage(buf.Bytes()).Size()
	if err != nil {
		t.Fatal(err)
	}
	if size.Width != 600 || size.Height != 800 {
		t.Fatalf("size = %d×%d, want 600×800", size.Width, size.Height)
	}
}
