package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var PORT string
var SERVICE_NAME string
var API_GATEWAY string
var CONTAINER_URL string

func main() {
	PORT = os.Getenv("PORT")
	SERVICE_NAME = os.Getenv("SERVICE_NAME")
	API_GATEWAY = os.Getenv("API_GATEWAY")
	CONTAINER_URL = os.Getenv("CONTAINER_URL")
	registerTheService()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World from service: %s <br> %s \n", SERVICE_NAME, r.URL.String())
	})

	http.HandleFunc("/all", allCommentsHandler)

	http.ListenAndServe(":"+PORT, nil)
}

func allCommentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All Comments goes here")
}

func registerTheService() {
	requestBody := map[string]string{
		"service_name": SERVICE_NAME,
		"url":          CONTAINER_URL,
	}

	body, _ := json.Marshal(requestBody)
	resp, err := http.Post(API_GATEWAY, "application/json", bytes.NewBuffer(body))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("Service has been registered")
}
