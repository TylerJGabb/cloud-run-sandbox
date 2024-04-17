package http_handlers

import (
	"cloud-run-sandbox/middleware"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewGetFileContentsHandler(filesLocation string) GetFileContentsHandler {
	return GetFileContentsHandler{
		filesLocation: filesLocation,
	}
}

type GetFileContentsHandler struct {
	filesLocation string
}


func (h GetFileContentsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := middleware.GetTraceLogger(*req)
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
	data, err := os.ReadFile(h.filesLocation + "/" + sd.Name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist)  {
			logger.Error(fmt.Sprintf("File not found: %v", err))
			w.WriteHeader(http.StatusNotFound)
		} else {
			logger.Error(fmt.Sprintf("Failed to read file: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		return
	}
	logger.Info(fmt.Sprintf("Contents: %v", string(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
