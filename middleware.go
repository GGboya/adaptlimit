package adaptlimit

import (
	"net/http"
)

// LimitHandler wraps an http.Handler with concurrency limiting middleware
func (l *Limiter) LimitHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release, err := l.Acquire(r.Context())
		if err != nil {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// 包装 ResponseWriter 以捕获状态码
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		defer func() {
			// 如果状态码小于 500，我们认为请求是成功的
			success := rw.status < 500
			release(success)
		}()

		next.ServeHTTP(rw, r)
	})
}

// LimitHandlerFunc wraps an http.HandlerFunc with concurrency limiting middleware
func (l *Limiter) LimitHandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.LimitHandler(next).ServeHTTP(w, r)
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
