package logging

import (
	c "cloud-run-sandbox/config"

	"context"
	"fmt"
	"net/http"
	"strings"
)

func newLoggerFromRequest(req http.Request) Logger {
	var trace string
	traceHeader := req.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = fmt.Sprintf("projects/%s/traces/%s", c.ProjectId, traceParts[0])
	}
	logger := Logger{
		trace: trace,
	}
	if trace != "" {
		logger.Info(fmt.Sprintf("Trace: %v", trace))
	} else {
		logger.Warn("Trace is not set")

	}
	return logger
}

type requestLoggerKey struct {}

func LoggerFromRequest(req *http.Request) Logger {
	return req.Context().Value(requestLoggerKey{}).(Logger)
}

func TraceMiddleware(handler func(w http.ResponseWriter, req *http.Request)) http.Handler {
	next := http.HandlerFunc(handler)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := newLoggerFromRequest(*req)
		ctx := context.WithValue(req.Context(), requestLoggerKey{}, logger)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}