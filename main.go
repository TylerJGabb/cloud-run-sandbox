package main

import (
	"cloud-run-sandbox/config"
	"cloud-run-sandbox/http_handlers"
	"cloud-run-sandbox/middleware"
	"cloud-run-sandbox/server"
	"log"
)

func init() {
	// TODO: just use fmt.Println, don't use log
	log.SetFlags(0)
}

func main() {
	// TODO: create a context and pass it down properly
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	app := server.NewAppServer()
	app.Use(middleware.WithTraceLogger(cfg.ProjectId))
	app.Use(middleware.SayHelloWithTraceLogger)
	app.Handle("/", http_handlers.NewGetFileContentsHandler(cfg.FilesLocation))
	app.Start(":" + cfg.Port)
}
