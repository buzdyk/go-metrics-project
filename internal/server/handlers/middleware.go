package handlers

import (
	"compress/gzip"
	"compress/zlib"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

type compressedResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (c *compressedResponseWriter) Write(b []byte) (int, error) {
	return c.Writer.Write(b)
}

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			logger.Info("request handled",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int("status", rw.statusCode),
				zap.Int("response_size", rw.size),
				zap.Duration("duration", duration),
			)
		})
	}
}

func DecompressRequestMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Encoding") == "gzip" {
				gzipReader, err := gzip.NewReader(r.Body)
				if err != nil {
					http.Error(w, "Invalid gzip data", http.StatusBadRequest)
					return
				}
				defer gzipReader.Close()
				r.Body = gzipReader
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CompressResponseMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encodings := r.Header.Get("Accept-Encoding")
			if strings.Contains(encodings, "gzip") {
				gzipWriter := gzip.NewWriter(w)
				defer gzipWriter.Close()
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Vary", "Accept-Encoding")
				next.ServeHTTP(&compressedResponseWriter{ResponseWriter: w, Writer: gzipWriter}, r)
				return
			}
			if strings.Contains(encodings, "deflate") {
				zlibWriter := zlib.NewWriter(w)
				defer zlibWriter.Close()
				w.Header().Set("Content-Encoding", "deflate")
				w.Header().Set("Vary", "Accept-Encoding")
				next.ServeHTTP(&compressedResponseWriter{ResponseWriter: w, Writer: zlibWriter}, r)
				return
			}
			next.ServeHTTP(w, r) // No compression if not supported
		})
	}
}
