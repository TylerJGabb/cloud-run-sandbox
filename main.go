package main

import (
	c "cloud-run-sandbox/config"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(0)
	if err := c.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		panic(err)
	}
	panic(http.ListenAndServe(":"+c.Port, nil))
}
