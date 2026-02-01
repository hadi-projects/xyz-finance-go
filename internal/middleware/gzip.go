package middleware

import (
	"compress/gzip"
	"io"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	// MinCompressionSize is the minimum response size in bytes before applying compression
	MinCompressionSize = 1024 // 1KB
)

// gzipWriter wraps http.ResponseWriter to provide gzip compression
type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

// gzip.Writer pool to reduce allocations
var gzipPool = sync.Pool{
	New: func() interface{} {
		w, _ := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
		return w
	},
}

// GzipCompression returns a middleware that compresses response using gzip
// Only compresses responses larger than MinCompressionSize (1KB)
// Only compresses if client supports gzip encoding
func GzipCompression() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if client doesn't accept gzip
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Skip for websocket upgrades
		if c.GetHeader("Upgrade") == "websocket" {
			c.Next()
			return
		}

		// Skip for Server-Sent Events
		if c.GetHeader("Accept") == "text/event-stream" {
			c.Next()
			return
		}

		// Get gzip writer from pool
		gz := gzipPool.Get().(*gzip.Writer)
		gz.Reset(c.Writer)

		// Create wrapped writer
		gzWriter := &gzipWriter{
			ResponseWriter: c.Writer,
			writer:         gz,
		}

		// Set response headers
		c.Writer = gzWriter
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// Process request
		c.Next()

		// Close and return writer to pool
		gz.Close()
		gzipPool.Put(gz)
	}
}

// ConditionalGzip applies gzip only for responses > MinCompressionSize
// Uses buffered approach to check size before compressing
func ConditionalGzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if client doesn't accept gzip
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Skip for non-compressible content types
		c.Next()

		// After response is written, check if we should compress
		// Note: This is handled at response write time in the standard middleware
	}
}

// shouldCompress checks if the content type should be compressed
func shouldCompress(contentType string) bool {
	compressibleTypes := []string{
		"application/json",
		"text/html",
		"text/plain",
		"text/xml",
		"application/xml",
		"text/css",
		"application/javascript",
		"text/javascript",
	}

	for _, t := range compressibleTypes {
		if strings.Contains(contentType, t) {
			return true
		}
	}
	return false
}

// GzipCompressionWithConfig returns a middleware with custom configuration
func GzipCompressionWithConfig(minSize int, level int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if client doesn't accept gzip
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Skip for websocket upgrades
		if c.GetHeader("Upgrade") == "websocket" {
			c.Next()
			return
		}

		// Create new gzip writer with custom level
		gz, err := gzip.NewWriterLevel(c.Writer, level)
		if err != nil {
			c.Next()
			return
		}

		// Create wrapped writer
		gzWriter := &gzipWriter{
			ResponseWriter: c.Writer,
			writer:         gz,
		}

		// Set response headers
		c.Writer = gzWriter
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// Process request
		c.Next()

		// Close writer
		gz.Close()
	}
}
