package main

import (
	"fmt"
	"net/http"
)

func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		allowed, _ := app.rateLimiter.Allow(ip)
		if !allowed {
			app.errorHandler.TooManyRequests(w, r, fmt.Errorf("too many requests"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
