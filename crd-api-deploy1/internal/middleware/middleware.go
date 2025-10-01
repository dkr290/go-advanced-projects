// Package middleware
package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Content-Type-Protection", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000;includeSubDomains;preload")
		w.Header().Set("Referred-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received Request in ResponceTime")
		start := time.Now()

		// creaste custom resp writer to capture the status code
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		duration := time.Since(start)
		wrappedWriter.Header().Set("X-Response-Time", duration.String())
		next.ServeHTTP(wrappedWriter, r)
		duration = time.Since(start)
		fmt.Printf(
			"Method: %s, URL: %s, Status: %d, Duration: %v \n",
			r.Method,
			r.URL,
			wrappedWriter.status,
			duration.String(),
		)
		fmt.Println("Send responce from Responce Time Middleware")
	})
}

// response writer
// create custom responce writer

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
