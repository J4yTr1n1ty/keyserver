package middleware

import (
	"log"
	"net/http"
	"time"
)

// We need to wrap the writer here to be able to extract the statuscode which was assigned by the handler
type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the start time of the request
		start := time.Now()

		// Create a new wrapped writer
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		// Use the wrapped writer in the next handler call.
		next.ServeHTTP(wrapped, r)

		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start)) // Log the request with the status code and elapsed time
	})
}
