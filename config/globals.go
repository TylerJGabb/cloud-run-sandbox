package config

import (
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
)


const (
	GOOGLE_PROJECT_ID_ENV = "GOOGLE_PROJECT_ID"
)

var (
	ProjectId   string
	Port string
)


func Load() error {
	if projectIdFromEnv, ok := os.LookupEnv(GOOGLE_PROJECT_ID_ENV); ok {
		ProjectId = projectIdFromEnv
	} else {
		// https://cloud.google.com/run/docs/container-contract#metadata-server
		// https://pkg.go.dev/cloud.google.com/go/compute/metadata
		log.Println(GOOGLE_PROJECT_ID_ENV + " is not set, trying to get project ID from metadata server")
		metadataClient := metadata.NewClient(&http.Client{
			Timeout: time.Duration(5 * time.Second),
		})
		metaProjectId, err := metadataClient.ProjectID()
		if err != nil {
			// errMsg := fmt.Sprintf("Failed to get project ID: %v", err)
			// RootLogger.Error(errMsg)
			return err
		}
		ProjectId = metaProjectId
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	Port = port
	// RootLogger.Info("Configuration Loaded Successfully", 
		// "ProjectID", ProjectId,
		// "Port", Port,
	// )
	return nil
}