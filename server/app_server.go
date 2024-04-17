package server

import (
	"cloud-run-sandbox/logging"
	"cloud-run-sandbox/middleware"
	"net/http"
)

func NewAppServer() *AppServer {
	return &AppServer{
		middlewares: make([]middleware.Middleware, 0),
		handles: make(map[string]http.Handler),
	}
}

type AppServer struct {
	middlewares []middleware.Middleware
	handles map[string]http.Handler
}

func (as *AppServer) Use(mid middleware.Middleware) {
	 as.middlewares = append(as.middlewares, mid)
}

func (as AppServer) Handle(pattern string, handler http.Handler) {
	as.handles[pattern] = handler
}

func (as AppServer) Start(addr string) {
	mux := http.NewServeMux()
	// Why do we call in reverse?
	// When someone adds handlers, they do so in the intuitive order
	// use(InjectLogger)
	// use(UseLogger)
	// but what this would actually look like in-line would be
	// finalHandler := InjectLogger(next=UseLogger(next=theHandler))
	chain := func(next http.Handler) http.Handler {
		for i := len(as.middlewares) - 1; i >= 0; i-- {
			mid := as.middlewares[i]
			next = mid(next)
		}
		return next
	}
	for pattern, handler := range as.handles {
		logging.SharedLogger.Debug("Adding Handler for " + pattern)
		withMiddleware := chain(handler)
		mux.Handle(pattern, withMiddleware)
	}
	logging.SharedLogger.Info("Starting server", "address", addr)
	http.ListenAndServe(addr, mux)
}
