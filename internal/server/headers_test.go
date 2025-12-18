package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeaders_RemovesHopHeaders(t *testing.T) {
	// response with hop headers + other headers
	resp := &http.Response{
		Header: http.Header{
			"Content-Type":      {"application/json"},
			"X-Test":            {"a", "b"},
			"Connection":        {"keep-alive"},
			"Transfer-Encoding": {"chunked"},
			"Upgrade":           {"websocket"},
		},
	}

	rr := httptest.NewRecorder()
	CopyResponseHeaders(rr, resp)
	h := rr.Header()

	if got := h.Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}

	values := h.Values("X-Test")
	if len(values) != 2 || values[0] != "a" || values[1] != "b" {
		t.Fatalf("X-Test values = %#v, want [a b]", values)
	}
}
