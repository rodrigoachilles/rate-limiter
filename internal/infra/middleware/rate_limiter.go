package middleware

import (
	"github.com/rodrigoachilles/rate-limiter/internal/usecase/limiter"
	"net/http"
	"strings"

	"context"
)

func RateLimiter(l *limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := getKey(r)
			limit := l.IPLimit
			if token := r.Header.Get("API_KEY"); token != "" {
				key = token
				limit = l.TokenLimit
			}

			allowed, err := l.AllowRequest(context.Background(), key, limit)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if !allowed {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getKey(r *http.Request) string {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	return "rl:" + ip
}
