package middleware

import (
	"cloud-run-sandbox/config"
	"cloud-run-sandbox/logging"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type traceLoggerKey struct{}

func GetTraceLogger(req http.Request) logging.Logger {
	return req.Context().Value(traceLoggerKey{}).(logging.Logger)
}

func InjectLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {
		var trace string
		logger := logging.Logger{}
		traceHeader := req.Header.Get("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			// TODO: refactor config into context
			trace = fmt.Sprintf("projects/%s/traces/%s", config.ProjectId, traceParts[0])
			logger.Trace = trace
			logger.Info(fmt.Sprintf("Trace: %v", trace))
		} else {
			logger.Warn("Trace is not set")
		}
		ctx := context.WithValue(req.Context(), traceLoggerKey{}, logger)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}