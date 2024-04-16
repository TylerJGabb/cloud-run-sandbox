package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	GOOGLE_PROJECT_ID_ENV = "GOOGLE_PROJECT_ID"
)

var (
	projectId   string
	projectIdOk bool
	rootLogger  Logger = Logger{}
)

type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`
	Extras map[string]interface{} `json:"extra,omitempty"`
}

// ToJson renders an entry structure to the JSON format expected by Cloud Logging.
func (e Entry) ToJson() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

func KeyValuesToExtras(keyValue ...any) map[string]any {
	var key string
	result := map[string]any{}
	for idx, item := range keyValue {
		if idx%2 == 0 {
			key = item.(string)
		} else {
			result[key] = item
		}
	}
	return result
}

func (l Logger) Error(m string, keyValues ...interface{}) {
	e := Entry{
		Message:     m,
		Severity:    "ERROR",
		Trace:       l.Trace,
		Extras: KeyValuesToExtras(keyValues...),
	}
	log.Println(e.ToJson())
}

func (l Logger) Info(m string, keyValues ...interface{}) {
	e := Entry{
		Message:     m,
		Trace:       l.Trace,
		Extras: KeyValuesToExtras(keyValues...),
	}
	log.Println(e.ToJson())
}

func (l Logger) Warn(m string, keyValues ...interface{}) {
	e := Entry{
		Message:     m,
		Severity:    "WARNING",
		Trace:       l.Trace,
		Extras: KeyValuesToExtras(keyValues...),
	}
	log.Println(e.ToJson())
}

type Logger struct {
	Trace string
}

func NewLoggerFromRequest(req http.Request) Logger {
	var trace string
	if projectId != "" {
		traceHeader := req.Header.Get("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			trace = fmt.Sprintf("projects/%s/traces/%s", projectId, traceParts[0])
		}
	}
	logger := Logger{
		Trace: trace,
	}
	if trace != "" {
		logger.Info(fmt.Sprintf("Trace: %v", trace))
	} else {
		logger.Warn("Trace is not set")

	}
	return logger
}

func handler(w http.ResponseWriter, req *http.Request) {
	logger := NewLoggerFromRequest(*req)
	sd := storagedata.StorageObjectData{}
	bytes, err := io.ReadAll(req.Body)
	logger.Info(fmt.Sprintf("Received event: %v", string(bytes)))
	if err != nil {
		errMsg := fmt.Sprintf("Failed to read request body: %v", err)
		logger.Error(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return
	}
	defer req.Body.Close()
	if err := protojson.Unmarshal(bytes, &sd); err != nil {
		errMsg := fmt.Sprintf("Failed to unmarshal request body: %v", err)
		logger.Error(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return
	}

	// print the data
	logger.Info("Item Summary",
		"Bucket", sd.Bucket,
		"Name", sd.Name,
		"Id", sd.Id,
	)

	// print the contents of the file /gcs/${name}

	data, err := os.ReadFile("/gcs/" + sd.Name)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read file: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	logger.Info(fmt.Sprintf("Contents: %v", string(data)))
	w.WriteHeader(http.StatusOK)
}

func main() {
	log.SetFlags(0)
	// https://cloud.google.com/run/docs/container-contract#metadata-server
	// https://pkg.go.dev/cloud.google.com/go/compute/metadata
	projectId, projectIdOk = os.LookupEnv(GOOGLE_PROJECT_ID_ENV)
	if !projectIdOk {
		rootLogger.Warn(GOOGLE_PROJECT_ID_ENV + " is not set, trying to get project ID from metadata server")
		metadataClient := metadata.NewClient(&http.Client{
			Timeout: time.Duration(5 * time.Second),
		})
		metaProj, err := metadataClient.ProjectID()
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get project ID: %v", err)
			rootLogger.Error(errMsg)
			panic(err)
		}
		projectId = metaProj
	}
	rootLogger.Info(fmt.Sprintf("Project ID: %v", projectId))
	http.HandleFunc("/", handler)
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	rootLogger.Info("Server starting at port " + port)
	panic(http.ListenAndServe(":"+port, nil))
}
