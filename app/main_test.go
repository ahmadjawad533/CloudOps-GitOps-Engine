package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	healthzHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("healthzHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK\n"
	if rr.Body.String() != expected {
		t.Errorf("healthzHandler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestReadyHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rr := httptest.NewRecorder()

	readyHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("readyHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	rootHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("rootHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "operational") {
		t.Errorf("rootHandler body expected to contain 'operational', got: %s", body)
	}
}

func TestMetricsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	metricsHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("metricsHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "http_requests_total") {
		t.Errorf("metricsHandler body missing 'http_requests_total', got: %s", body)
	}
}
