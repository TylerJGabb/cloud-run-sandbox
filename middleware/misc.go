package middleware

import "net/http"

func SayHelloWithTraceLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {
		tl := GetTraceLogger(*req)
		tl.Info("Hello")
		next.ServeHTTP(w, req)
	})
}