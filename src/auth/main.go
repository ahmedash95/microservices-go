package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var PORT string
var SERVICE_NAME string
var API_GATEWAY string
var CONTAINER_URL string

var hmacSampleSecret = []byte("ABC123!@#")

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

func main() {
	PORT = os.Getenv("PORT")
	SERVICE_NAME = os.Getenv("SERVICE_NAME")
	API_GATEWAY = os.Getenv("API_GATEWAY")
	CONTAINER_URL = os.Getenv("CONTAINER_URL")

	DB_Init()
	registerTheService()

	r := mux.NewRouter()
	r.HandleFunc("/user/create", HandleUserCreate).Methods("POST")
	r.HandleFunc("/login", HandleUserLogin).Methods("POST")
	r.HandleFunc("/user", HandleGetLoggedInUser).Methods("POST")

	http.ListenAndServe(":"+PORT, r)
}
