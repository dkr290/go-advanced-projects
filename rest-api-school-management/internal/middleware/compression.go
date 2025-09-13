package middleware

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if the clienbt accepts gzip encoding

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
		}

		// if it accepts gzip then set the responce header
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// wrap the gzipWriter now
		gzw := &gzipWriter{
			ResponseWriter: w,
			gz:             gz,
		}

		next.ServeHTTP(gzw, r)
		// logging alwqays after the next funxtion
		fmt.Println("Send responce from Compression Middleware")
	})
}

// we need to copmpress the responce
// gzip responce writer with wraps http.ResponceWriter to write gzip responces
// To properly compress the response, you need to wrap the `http.ResponseWriter` so that all writes go through the gzip writer. This requires a new type (e.g., `gzipResponseWriter`) that implements `http.ResponseWriter` and writes to the gzip stream.

type gzipWriter struct {
	http.ResponseWriter
	gz *gzip.Writer
}

func (g *gzipWriter) Write(p []byte) (int, error) {
	return g.gz.Write(p)
}
