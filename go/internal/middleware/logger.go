package middleware

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime"
	"time"
)

// Logger logs incoming request in a standardized format.
//
//	method=POST url=/users.get status=200 size=42 duration=10ms
//
// The middleware inject a [LogFields] object into the context allowing the next
// handlers to add custom fields to the log line:
//
//	method=POST url=/users.get status=200 size=42 duration=10ms team=xxxx-xxxx-xxxx-xxxxxxxx
func Logger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		runtime.GC()

		start := time.Now()
		wlog := newResponseLogger(w)

		next.ServeHTTP(wlog, r)
		if wlog.status == 0 {
			wlog.status = http.StatusOK
		}

		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)

		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.Int("status", wlog.status),
			slog.String("size", formatByteCount(uint64(wlog.size))),
			slog.String("heap", formatByteCount(stats.HeapAlloc)),
			slog.Duration("duration", time.Since(start).Round(time.Millisecond)),
		}
		logger.LogAttrs(context.Background(), slog.LevelInfo, "incoming request", attrs...)
	})
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

// newResponseLogger creates a new responseLogger or return an existing one.
func newResponseLogger(w http.ResponseWriter) *responseLogger {
	if tmp, ok := w.(*responseLogger); ok {
		return tmp
	}
	return &responseLogger{w: w}
}

// Header returns logger header.
func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

// Write implements [io.Writer].
func (l *responseLogger) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

// ReadFrom speeds up writing of large object.
func (l *responseLogger) ReadFrom(src io.Reader) (n int64, err error) {
	size, err := io.Copy(l.w, src)
	l.size += int(size)
	return size, err
}

// WriteHeader writes HTTP header with the given code.
func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

// Status returns the status written to the user.
func (l *responseLogger) Status() int {
	return l.status
}

// Size returns the amount of data sent to the user.
func (l *responseLogger) Size() int {
	return l.size
}

// Flush empty the response writer.
func (l *responseLogger) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

func formatByteCount(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "kMGTPE"[exp])
}
