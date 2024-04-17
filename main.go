package main

import (
	c "cloud-run-sandbox/config"
	"cloud-run-sandbox/http_handlers"
	"cloud-run-sandbox/middleware"
	"cloud-run-sandbox/server"
	"log"
	"net/http"
)

func main() {
	// TODO: create a context and pass it down properly
	log.SetFlags(0)
	if err := c.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		panic(err)
	}
	// TODO: abstract initializer that can take any App interface
	app := server.NewAppServer()
	app.Use(middleware.InjectLogger)
	app.Use(middleware.SayHelloWithLogger)
	app.Handle("/", http.HandlerFunc(http_handlers.GetFileContents))
	app.Start(":8080")
}
