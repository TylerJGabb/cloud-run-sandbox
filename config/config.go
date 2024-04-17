package config

import (
	"cloud-run-sandbox/logging"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
)

type Config struct {
	ProjectId string
	Port string
	FilesLocation string
}

const (
	filesLocationEnvVar = "FILES_LOCATION"
	portEnvVar = "PORT"
	googleProjectIdEnvVar = "GOOGLE_PROJECT_ID"
)

func loadFilesLocation() string {
	filesLocation, ok := os.LookupEnv(filesLocationEnvVar)
	if !ok {
		logging.SharedLogger.Warn(filesLocationEnvVar + " was not set, defaulting to '/gcs'")
		filesLocation = "/gcs"
	}
	return filesLocation
}

func loadPort() string {
	port, ok := os.LookupEnv(portEnvVar)
	if !ok {
		logging.SharedLogger.Warn(portEnvVar + " was not set, defaulting to 8080")
		port = "8080"
	}
	return port
}

func loadProjectId() (string, error) {
	if projectIdFromEnv, ok := os.LookupEnv(googleProjectIdEnvVar); ok {
		return projectIdFromEnv, nil
	} else {
		// https://cloud.google.com/run/docs/container-contract#metadata-server
		// https://pkg.go.dev/cloud.google.com/go/compute/metadata
		log.Println(googleProjectIdEnvVar + " is not set, trying to get project ID from metadata server")
		metadataClient := metadata.NewClient(&http.Client{
			Timeout: time.Duration(1 * time.Second),
		})
		metaProjectId, err := metadataClient.ProjectID()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get project ID: %v", err)
			logging.SharedLogger.Error(errMsg)
			return "", err
		}
		return metaProjectId, nil
	}
}

func Load() (Config, error) {
	projectId, err := loadProjectId()
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		FilesLocation: loadFilesLocation(),
		Port: loadPort(),
		ProjectId: projectId,
	}
	logging.SharedLogger.Info("Configuration Loaded Successfully", 
		"config", cfg,
	)
	return cfg, nil
}