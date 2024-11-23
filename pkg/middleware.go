package pkg

// import (
// 	"compress/gzip"
// 	"net/http"
// 	"strings"
// )

// // gzipMiddleware compresses HTTP responses using GZIP.
// func gzipMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Check if the client supports gzip encoding
// 		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		// Wrap the response writer with a GZIP writer
// 		w.Header().Set("Content-Encoding", "gzip")
// 		gz := gzip.NewWriter(w)
// 		defer gz.Close()

// 		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
// 		next.ServeHTTP(gzr, r)
// 	})
// }

// // gzipResponseWriter wraps http.ResponseWriter and gzip.Writer
// type gzipResponseWriter struct {
// 	http.ResponseWriter
// 	Writer *gzip.Writer
// }

// func (gz gzipResponseWriter) Write(b []byte) (int, error) {
// 	return gz.Writer.Write(b)
// }

// // customHeadersMiddleware adds custom headers to HTTP responses.
// func customHeadersMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Add custom headers
// 		w.Header().Set("X-Content-Type-Options", "nosniff")
// 		w.Header().Set("X-Frame-Options", "DENY")
// 		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
// 		w.Header().Set("Content-Security-Policy", "default-src 'self'")

// 		next.ServeHTTP(w, r)
// 	})
// }
