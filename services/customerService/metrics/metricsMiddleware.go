package metrics

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Ensure headers like CORS are not affected
		if rw.Header().Get("Access-Control-Allow-Origin") == "" {
			log.Println("CORS headers missing after MetricsMiddleware")
		}

		// Record metrics
		duration := time.Since(start).Seconds()
		HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rw.status)).Inc()
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
