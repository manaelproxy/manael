// Copyright (c) 2024 Yamagishi Kazutoshi <ykzts@desire.sh>
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

//go:build integration

package manael_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

// webpRIFF is the 4-byte RIFF container signature at the start of every WebP file.
var webpRIFF = []byte("RIFF")

// webpFourCC is the 4-byte "WEBP" identifier found at byte offset 8 in every WebP file.
var webpFourCC = []byte("WEBP")

// isWebP returns true if b begins with a valid WebP header: RIFF + 4-byte size + WEBP.
func isWebP(b []byte) bool {
	return len(b) >= 12 &&
		bytes.Equal(b[:4], webpRIFF) &&
		bytes.Equal(b[8:12], webpFourCC)
}

const (
	nginxAlias  = "upstream"
	manaelPort  = "8080/tcp"
	nginxPort   = "80/tcp"
)

// setupContainers starts an Nginx container (upstream) and a Manael container
// (proxy) connected to a shared Docker network, and returns the Manael base URL.
// The returned cleanup function stops and removes both containers and the network.
func setupContainers(ctx context.Context, t *testing.T) string {
	t.Helper()

	repoRoot, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	// Create an isolated Docker network for container-to-container communication.
	net, err := network.New(ctx)
	if err != nil {
		t.Fatalf("create network: %v", err)
	}
	t.Cleanup(func() {
		if err := net.Remove(ctx); err != nil {
			t.Logf("remove network: %v", err)
		}
	})

	networkName := net.Name

	// Start Nginx container to serve testdata/ as the upstream image origin.
	nginxReq := testcontainers.ContainerRequest{
		Image: "nginx:alpine",
		Mounts: testcontainers.Mounts(
			testcontainers.BindMount(
				filepath.Join(repoRoot, "testdata"),
				"/usr/share/nginx/html",
			),
		),
		Networks: []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {nginxAlias},
		},
		ExposedPorts: []string{nginxPort},
		WaitingFor:   wait.ForHTTP("/logo.png").WithPort(nginxPort),
	}
	nginxContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: nginxReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start nginx container: %v", err)
	}
	t.Cleanup(func() {
		if err := nginxContainer.Terminate(ctx); err != nil {
			t.Logf("terminate nginx container: %v", err)
		}
	})

	// Build and start the Manael container from the local Dockerfile.
	manaelReq := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: repoRoot,
		},
		Env: map[string]string{
			"MANAEL_UPSTREAM_URL": fmt.Sprintf("http://%s:80", nginxAlias),
		},
		Networks:     []string{networkName},
		ExposedPorts: []string{manaelPort},
		WaitingFor:   wait.ForListeningPort(manaelPort),
	}
	manaelContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: manaelReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start manael container: %v", err)
	}
	t.Cleanup(func() {
		if err := manaelContainer.Terminate(ctx); err != nil {
			t.Logf("terminate manael container: %v", err)
		}
	})

	host, err := manaelContainer.Host(ctx)
	if err != nil {
		t.Fatalf("get manael host: %v", err)
	}
	port, err := manaelContainer.MappedPort(ctx, manaelPort)
	if err != nil {
		t.Fatalf("get manael port: %v", err)
	}

	return fmt.Sprintf("http://%s:%s", host, port.Port())
}

// TestIntegration runs the full end-to-end integration suite against a live
// Manael Docker container proxying an Nginx upstream.
func TestIntegration(t *testing.T) {
	ctx := context.Background()
	baseURL := setupContainers(ctx, t)

	t.Run("PNGToWebP", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/logo.png", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Accept", "image/webp,image/*,*/*;q=0.8")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("status code = %d, want %d", got, want)
		}
		if got, want := resp.Header.Get("Content-Type"), "image/webp"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !isWebP(body) {
			t.Errorf("response body does not contain valid WebP RIFF magic bytes")
		}
	})

	t.Run("JPEGToWebP", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/photo.jpeg", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Accept", "image/webp,image/*,*/*;q=0.8")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("status code = %d, want %d", got, want)
		}
		if got, want := resp.Header.Get("Content-Type"), "image/webp"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !isWebP(body) {
			t.Errorf("response body does not contain valid WebP RIFF magic bytes")
		}
	})

	t.Run("NoWebPAcceptPassthrough", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/logo.png", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Accept", "image/png,image/*,*/*;q=0.8")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("status code = %d, want %d", got, want)
		}
		if got, want := resp.Header.Get("Content-Type"), "image/png"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}
	})

	t.Run("NotFoundPassthrough", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/nonexistent.png", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Accept", "image/webp,image/*,*/*;q=0.8")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if got, want := resp.StatusCode, http.StatusNotFound; got != want {
			t.Errorf("status code = %d, want %d", got, want)
		}
	})

	t.Run("NonImagePassthrough", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/empty.txt", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Accept", "image/webp,image/*,*/*;q=0.8")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("status code = %d, want %d", got, want)
		}
		ct := resp.Header.Get("Content-Type")
		if ct == "image/webp" || ct == "image/avif" {
			t.Errorf("Content-Type = %q, want non-image type (text file should pass through)", ct)
		}
	})

	t.Run("VaryHeaderPresent", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/logo.png", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Accept", "image/webp,image/*,*/*;q=0.8")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		vary := resp.Header.Get("Vary")
		if vary == "" {
			t.Error("Vary header is absent, want it to contain \"Accept\"")
		}
	})

	t.Run("ServerHeader", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/logo.png", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if got, want := resp.Header.Get("Server"), "Manael"; got != want {
			t.Errorf("Server = %q, want %q", got, want)
		}
	})
}
