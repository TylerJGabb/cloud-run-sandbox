package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
)

/*
TODO: build some sort of logger that extracts the tract id

	Derive the traceID associated with the current request.
	var trace string
	if projectID != "" {
	        traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	        traceParts := strings.Split(traceHeader, "/")
	        if len(traceParts) > 0 && len(traceParts[0]) > 0 {
	                trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
	        }
	}
*/
type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`

	// Logs Explorer allows filtering and display of this as `jsonPayload.component`.
	Component string `json:"component,omitempty"`
}

// String renders an entry structure to the JSON format expected by Cloud Logging.
func (e Entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

func Error(m string) {
	e := Entry{Message: m, Severity: "ERROR"}
	log.Println(e.String())
}

func Info(m string) {
	e := Entry{Message: m}
	log.Println(e.String())
}

func Warn(m string) {
	e := Entry{Message: m, Severity: "WARNING"}
	log.Println(e.String())
}

func handler(w http.ResponseWriter, req *http.Request) {
	sd := storagedata.StorageObjectData{}
	bytes, err := io.ReadAll(req.Body)
	Info(fmt.Sprintf("Received event: %v\n", string(bytes)))
	if err != nil {
		errMsg := fmt.Sprintf("Failed to read request body: %v", err)
		Error(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return
	}
	defer req.Body.Close()
	if err := protojson.Unmarshal(bytes, &sd); err != nil {
		errMsg := fmt.Sprintf("Failed to unmarshal request body: %v", err)
		Error(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return

	}

	log.Println()

	// print the data
	fmt.Printf("Bucket: %v\n", sd.Bucket)
	fmt.Printf("Name: %v\n", sd.Name)
	fmt.Printf("Id: %v\n", sd.Id)

	// print the contents of the file /gcs/${name}

	data, err := os.ReadFile("/gcs/" + sd.Name)
	if err != nil {
		Error(fmt.Sprintf("Failed to read file: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	Info(fmt.Sprintf("Contents: %v\n", string(data)))
	w.WriteHeader(http.StatusOK)
}

func main() {
	log.SetFlags(0)
	http.HandleFunc("/", handler)
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	fmt.Printf("Server starting at port %s\n", port)
	panic(http.ListenAndServe(":"+port, nil))
}
