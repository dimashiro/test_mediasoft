package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logging(log *zap.SugaredLogger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Infow("request completed", "method", r.Method, "path", r.URL.Path,
			"remoteaddr", r.RemoteAddr, "since", time.Since(start))
	}
}
