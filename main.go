package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
)


func handler(w http.ResponseWriter, req *http.Request) {
	sd := storagedata.StorageObjectData{}
	bytes, err := ioutil.ReadAll(req.Body)
	fmt.Printf("Received event: %v\n", string(bytes))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer req.Body.Close()
	if err := protojson.Unmarshal(bytes, &sd); err != nil {
		fmt.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	
	}
	// print the data
	fmt.Printf("Bucket: %v\n", sd.Bucket)
	fmt.Printf("Name: %v\n", sd.Name)
	fmt.Printf("Id: %v\n", sd.Id)
	w.WriteHeader(http.StatusOK)
	
}

func main() {
	http.HandleFunc("/", handler)
	print("Server starting at port 8080")
	panic(http.ListenAndServe(":8080", nil))
}

