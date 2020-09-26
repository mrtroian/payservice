package api

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var serverName string

func timer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		finish := time.Since(start)
		responseTime := finish.Microseconds()
		processedIn := fmt.Sprintf("%d μs", responseTime)

		w.Header().Add("X-Response-Time", processedIn)
		log.Println(r.Method, r.URL.Path, processedIn)
		return
	})
}

func get(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
