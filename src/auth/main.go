package main

import (
	"net/http"
	"os"

	"github.com/ahmedash95/gatewaySDK"
	"github.com/gorilla/mux"
)

var PORT string
var SERVICE_NAME string
var CONTAINER_URL string

var hmacSampleSecret = []byte("ABC123!@#")

func main() {
	PORT = os.Getenv("PORT")
	SERVICE_NAME = os.Getenv("SERVICE_NAME")
	CONTAINER_URL = os.Getenv("CONTAINER_URL")

	service := gatewaySDK.Service{
		SERVICE_NAME,
		CONTAINER_URL,
	}
	_, err := gatewaySDK.RegisterService(service)
	if err != nil {
		panic(err)
	}

	DB_Init()

	r := mux.NewRouter()
	r.HandleFunc("/user/create", HandleUserCreate).Methods("POST")
	r.HandleFunc("/login", HandleUserLogin).Methods("POST")
	r.HandleFunc("/user", HandleGetLoggedInUser).Methods("POST")

	http.ListenAndServe(":"+PORT, r)
}
