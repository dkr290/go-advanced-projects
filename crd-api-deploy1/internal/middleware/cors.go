package middleware

import (
	"net/http"
	"slices"
)

var allowedOrigins = []string{
	"https://myorigin.example.com",
	"http://localhost:8080",
	"http://localhost:8082",
	"http://localhost:8089",
	"https://k8s-dev.domain.com",
	"https://github.com",
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// fmt.Printf("Origin header: %q\n", origin)
		// fmt.Printf("Allowed origins: %v\n", allowedOrigins)
		// fmt.Printf("Is allowed: %v\n", isOriginAllowed(origin))
		// this allows when client does not set origin , not good for security but some browsers dont set this so just for testing to remove in production
		if origin != "" {
			if !isOriginAllowed(origin) {
				http.Error(w, "Not allowed by CORS", http.StatusForbidden)
				return
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string) bool {
	return slices.Contains(allowedOrigins, origin)
}
