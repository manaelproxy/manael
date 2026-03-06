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
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/noop"
	"manael.org/x/manael/v2"
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

	otel.SetMeterProvider(noop.NewMeterProvider())

	var handler http.Handler
	handler = manael.NewServeProxy(upstreamURL)
	handler = otelhttp.NewHandler(handler, "manael-proxy")
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
