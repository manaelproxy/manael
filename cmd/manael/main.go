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

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric"
	"manael.org/x/manael/v3"
)

const (
	// DefaultPort is returned by default port.
	DefaultPort = 8080

	// DefaultUpstreamURL is returned by default upstream URL.
	DefaultUpstreamURL = "http://localhost:9000"
)

type config struct {
	httpAddr    string
	upstreamURL string
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	conf := config{}

	fs.StringVar(&conf.httpAddr, "http", "", "HTTP server address")
	fs.StringVar(&conf.upstreamURL, "upstream_url", "", "Upstream URL for processing images")

	if err := fs.Parse(os.Args[1:]); err != nil {
		slog.Error("failed to parse flags", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if conf.httpAddr == "" {
		port := os.Getenv("PORT")

		if port != "" {
			conf.httpAddr = fmt.Sprintf(":%s", port)
		} else {
			conf.httpAddr = fmt.Sprintf(":%d", DefaultPort)
		}
	}

	if conf.upstreamURL == "" {
		u := os.Getenv("MANAEL_UPSTREAM_URL")

		if u != "" {
			conf.upstreamURL = u
		} else {
			conf.upstreamURL = DefaultUpstreamURL
		}
	}

	upstreamURL, err := url.Parse(conf.upstreamURL)
	if err != nil {
		slog.Error("failed to parse upstream URL", slog.String("error", err.Error()))
		os.Exit(1)
	}

	enableAVIF := os.Getenv("MANAEL_ENABLE_AVIF") == "true"
	enableResize := os.Getenv("MANAEL_ENABLE_RESIZE") == "true"

	var proxyOpts []manael.ProxyOption
	proxyOpts = append(proxyOpts, manael.WithAVIFEnabled(enableAVIF))
	proxyOpts = append(proxyOpts, manael.WithResizeEnabled(enableResize))

	if s := os.Getenv("MANAEL_MAX_IMAGE_SIZE"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil && n > 0 {
			proxyOpts = append(proxyOpts, manael.WithMaxImageSize(n))
		}
	}

	if s := os.Getenv("MANAEL_MAX_RESIZE_WIDTH"); s != "" {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || n <= 0 {
			slog.Error("invalid MANAEL_MAX_RESIZE_WIDTH",
				slog.String("error", fmt.Sprintf("must be a positive integer, got %q", s)))
			os.Exit(1)
		}
		proxyOpts = append(proxyOpts, manael.WithMaxResizeWidth(n))
	}

	if s := os.Getenv("MANAEL_MAX_RESIZE_HEIGHT"); s != "" {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || n <= 0 {
			slog.Error("invalid MANAEL_MAX_RESIZE_HEIGHT",
				slog.String("error", fmt.Sprintf("must be a positive integer, got %q", s)))
			os.Exit(1)
		}
		proxyOpts = append(proxyOpts, manael.WithMaxResizeHeight(n))
	}

	if s := os.Getenv("MANAEL_ALLOWED_WIDTHS"); s != "" {
		widths, err := parseIntList(s)
		if err != nil {
			slog.Error("invalid MANAEL_ALLOWED_WIDTHS", slog.String("error", err.Error()))
			os.Exit(1)
		}
		if len(widths) > 0 {
			proxyOpts = append(proxyOpts, manael.WithAllowedWidths(widths))
		}
	}

	if s := os.Getenv("MANAEL_ALLOWED_HEIGHTS"); s != "" {
		heights, err := parseIntList(s)
		if err != nil {
			slog.Error("invalid MANAEL_ALLOWED_HEIGHTS", slog.String("error", err.Error()))
			os.Exit(1)
		}
		if len(heights) > 0 {
			proxyOpts = append(proxyOpts, manael.WithAllowedHeights(heights))
		}
	}

	if s := os.Getenv("MANAEL_DEFAULT_QUALITY"); s != "" {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || n <= 0 {
			slog.Error("invalid MANAEL_DEFAULT_QUALITY",
				slog.String("error", fmt.Sprintf("must be a positive integer, got %q", s)))
			os.Exit(1)
		}
		proxyOpts = append(proxyOpts, manael.WithDefaultQuality(n))
	}

	if metricsPort := os.Getenv("MANAEL_METRICS_PORT"); metricsPort != "" {
		exporter, err := prometheus.New()
		if err != nil {
			slog.Error("failed to initialize prometheus exporter", slog.String("error", err.Error()))
			os.Exit(1)
		}
		provider := metric.NewMeterProvider(metric.WithReader(exporter))
		otel.SetMeterProvider(provider)

		metricsMux := http.NewServeMux()
		metricsMux.Handle("/metrics", promhttp.Handler())
		metricsAddr := ":" + metricsPort
		metricsSrv := &http.Server{
			Addr:              metricsAddr,
			Handler:           metricsMux,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
		}
		slog.Info("Starting metrics server", slog.String("addr", metricsAddr))
		go func() {
			if err := metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("metrics server error", slog.String("addr", metricsAddr), slog.String("error", err.Error()))
			}
		}()
	} else {
		otel.SetMeterProvider(noop.NewMeterProvider())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/_health", healthHandler)

	proxyHandler := manael.NewServeProxy(upstreamURL, proxyOpts...)
	mux.Handle("/", otelhttp.NewHandler(proxyHandler, "manael-proxy"))

	var handler http.Handler
	handler = mux
	handler = handlers.CombinedLoggingHandler(os.Stdout, handler)

	srv := &http.Server{
		Addr:              conf.httpAddr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("Starting server", slog.String("addr", conf.httpAddr))
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	stop()
	slog.Info("Shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(timeoutCtx); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Server exiting")
}

// healthHandler responds with 200 OK and "OK" for use as a liveness/readiness probe.
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// parseIntList splits a comma-separated string of positive integers. It
// returns an error for any token that is not a positive integer so that
// misconfiguration is caught at startup rather than silently disabling
// the whitelist safety controls.
func parseIntList(s string) ([]int, error) {
	var result []int
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		n, err := strconv.Atoi(part)
		if err != nil || n <= 0 {
			return nil, fmt.Errorf("must be a positive integer, got %q", part)
		}
		result = append(result, n)
	}
	return result, nil
}
