package server

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/hrz8/altalune"
)

func (s *Server) setupMiddleware(handler http.Handler) http.Handler {
	// apply middleware in reverse order (last applied executes first)
	handler = RecoveryMiddleware(handler, s.log)
	if s.cfg.IsHTTPLoggingEnabled() {
		handler = LoggingMiddleware(handler, s.log)
	}
	handler = SecurityMiddleware(handler)
	if s.cfg.IsCORSEnabled() {
		handler = CORSMiddleware(handler, s.cfg.GetAllowedOrigins())
	}

	return handler
}

func CORSMiddleware(next http.Handler, allowedOrigins []string) http.Handler {
	allowAll := false
	originMap := make(map[string]bool)
	for _, origin := range allowedOrigins {
		if origin == "*" {
			allowAll = true
			break
		}
		originMap[origin] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if allowAll || originMap[origin] {
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if allowAll {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Connect-Protocol-Version, Connect-Timeout-Ms")
			w.Header().Set("Access-Control-Expose-Headers", "Connect-Protocol-Version, Connect-Timeout-Ms")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler, log altalune.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		log.Info("incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.Header.Get("User-Agent"),
			"content_type", r.Header.Get("Content-Type"),
		)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", rw.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

func RecoveryMiddleware(next http.Handler, log altalune.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic recovered",
					"error", err,
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
					"stack", string(debug.Stack()),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal Server Error", "message": "An unexpected error occurred"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
