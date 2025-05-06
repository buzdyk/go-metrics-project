package handlers

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/sha256"
	"encoding/hex"
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

type hashResponseWriter struct {
	http.ResponseWriter
	key string
}

func (h *hashResponseWriter) Write(b []byte) (int, error) {
	// If a key is provided, calculate and set the hash header
	if h.key != "" {
		hash := sha256.New()
		hash.Write([]byte(string(b) + h.key))
		hashValue := hex.EncodeToString(hash.Sum(nil))
		h.ResponseWriter.Header().Set("HashSHA256", hashValue)
	}
	
	return h.ResponseWriter.Write(b)
}

type bodyReader struct {
	reader io.ReadCloser
	buffer *bytes.Buffer
}

func newBodyReader(r io.ReadCloser) *bodyReader {
	b := new(bytes.Buffer)
	return &bodyReader{
		reader: r,
		buffer: b,
	}
}

func (b *bodyReader) Read(p []byte) (int, error) {
	n, err := b.reader.Read(p)
	if n > 0 {
		b.buffer.Write(p[:n])
	}
	return n, err
}

func (b *bodyReader) Close() error {
	return b.reader.Close()
}

func (b *bodyReader) getBuffer() *bytes.Buffer {
	return b.buffer
}

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			logger.Info("request handled",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int("status", rw.statusCode),
				zap.Int("response_size", rw.size),
				zap.Duration("duration", duration),
				zap.Any("body", r.Body),
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

func ResponseHashMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wrap the response writer with our hash writer
			hashWriter := &hashResponseWriter{
				ResponseWriter: w,
				key:            key,
			}
			
			// Call the next handler with the wrapped writer
			next.ServeHTTP(hashWriter, r)
		})
	}
}

func SignatureVerificationMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If no key is provided, skip verification
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			
			// Check if the request has the HashSHA256 header
			hash := r.Header.Get("HashSHA256")
			if hash == "" {
				next.ServeHTTP(w, r)
				return
			}
			
			// For URL path parameters, extract the value from the URL
			var value string
			parts := strings.Split(r.URL.Path, "/")
			
			// Check if it's an update request with URL parameters
			if len(parts) >= 5 && parts[1] == "update" {
				// Value is the last part of the URL
				value = parts[len(parts)-1]
			} else {
				// For requests with a body (like JSON updates), read the body
				bodyReader := newBodyReader(r.Body)
				r.Body = bodyReader
				
				// Read the body
				body, err := io.ReadAll(bodyReader)
				if err != nil {
					http.Error(w, "Error reading request body", http.StatusBadRequest)
					return
				}
				
				// Reset the body for handlers
				r.Body = io.NopCloser(bytes.NewBuffer(body))
				value = string(body)
			}
			
			// Calculate the expected hash
			h := sha256.New()
			h.Write([]byte(value + key))
			expectedHash := hex.EncodeToString(h.Sum(nil))
			
			// Verify the hash
			if hash != expectedHash {
				http.Error(w, "Invalid signature", http.StatusBadRequest)
				return
			}
			
			// Continue to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
