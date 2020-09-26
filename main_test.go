package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_healthCheckHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(healthCheckHandler)
	handler.ServeHTTP(rr, req)

	want := http.StatusOK
	if got := rr.Code; got != want {
		t.Errorf("response code: got %v, want %v", got, want)
	}

	if got := rr.Body.String(); got != "ok" {
		t.Errorf("response code: got %v, want %v", got, want)
	}
}
