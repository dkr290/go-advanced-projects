// Package middleware
package middleware

import (
	"net/http"
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
