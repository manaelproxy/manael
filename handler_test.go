package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var us *httptest.Server

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	mux.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("200"))
		w.WriteHeader(200)
	})
	us = httptest.NewServer(mux)
	defer us.Close()

	ret := m.Run()

	us = nil
	os.Exit(ret)
}

func TestServeHTTP(t *testing.T) {
	h := &Handler{
		upstreamURL: us.URL,
	}
	ts := httptest.NewServer(h)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/200")
	if err != nil {
		t.Error("Could not connect.")
		return
	}

	if resp.StatusCode != 200 {
		t.Error("Status code incorrect")
		return
	}
}
