package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.NewString()
		}

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		ctx := r.Context()
		ctx = contextWithRequestInfo(ctx, requestID, traceID)
		r = r.WithContext(ctx)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		log.Printf("[%s] %s %s %d %s (reqID=%s traceID=%s)",
			r.Method,
			r.RemoteAddr,
			r.URL.Path,
			lrw.statusCode,
			duration,
			requestID,
			traceID,
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

type contextKey string

const (
	requestIDKey contextKey = "requestID"
	traceIDKey   contextKey = "traceID"
)

func contextWithRequestInfo(ctx context.Context, requestID, traceID string) context.Context {
	ctx = context.WithValue(ctx, requestIDKey, requestID)
	ctx = context.WithValue(ctx, traceIDKey, traceID)
	return ctx
}
