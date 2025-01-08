package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Info("Incoming Request",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remoteAddr", r.RemoteAddr),
				slog.String("userAgent", r.UserAgent()),
				slog.String("params", r.URL.RawQuery),
			)
			statusWriter := &StatusRecorder{ResponseWriter: w, StatusCode: http.StatusOK}
			next.ServeHTTP(statusWriter, r)
			duration := time.Since(start)

			if statusWriter.StatusCode >= 400 {
				responseBody := statusWriter.Body.String()
				var responseMap map[string]interface{}

				if err := json.Unmarshal([]byte(responseBody), &responseMap); err != nil {
					responseMap = map[string]interface{}{"rawResponse": responseBody}
				}

				logger.Error("Error",
					slog.String("requestURI", r.RequestURI),
					slog.Int("statusCode", statusWriter.StatusCode),
					slog.String("method", r.Method),
					slog.Any("responseBody", responseMap),
					slog.Duration("duration", duration),
				)
			} else {
				logger.Info("Request Handled",
					slog.String("method", r.Method),
					slog.String("url", r.URL.String()),
					slog.Int("statusCode", statusWriter.StatusCode),
					slog.Duration("duration", duration),
				)
			}
		})
	}
}

type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       bytes.Buffer
}

func (sr *StatusRecorder) WriteHeader(code int) {
	sr.StatusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *StatusRecorder) Write(b []byte) (int, error) {
	n, err := sr.Body.Write(b)
	if err != nil {
		return n, err
	}
	return sr.ResponseWriter.Write(b)
}
