package main

import (
	"cloud-run-sandbox/config"
	"cloud-run-sandbox/http_handlers"
	"cloud-run-sandbox/logging"
	"cloud-run-sandbox/middleware"
	"cloud-run-sandbox/server"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/google/uuid"
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
	stack := string(debug.Stack())
	logging.SharedLogger.Info("this is a trace:\n\n"+stack, "stack", stack, "foo", "bar", "number", 1234)
	app := server.NewAppServer()
	app.Use(middleware.WithTraceLogger(cfg.ProjectId))
	app.Use(middleware.SayHelloWithTraceLogger)
	app.Handle("/", http_handlers.NewGetFileContentsHandler(cfg.FilesLocation))
	app.Handle("/createFile", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := uuid.New().String()
		fName := cfg.FilesLocation + "/" + uid
		os.WriteFile(fName, []byte("Hello, World!"), 0644)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fName))
	}))
	app.Start(":" + cfg.Port)
}
