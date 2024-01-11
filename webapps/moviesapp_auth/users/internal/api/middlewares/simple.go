package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Simple(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(time.Since(start))
	})
}
