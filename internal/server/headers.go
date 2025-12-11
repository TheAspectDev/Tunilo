package server

import (
	"net/http"
)

// forbidden to forward ( src: RFC 7230 )
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

// sanitizes and copies headers from the connection response to the writer.
func CopyResponseHeaders(w http.ResponseWriter, resp *http.Response) {
	for _, h := range hopHeaders {
		resp.Header.Del(h)
	}

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	w.Header().Set("Connection", "close")
}
