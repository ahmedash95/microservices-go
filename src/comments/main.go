package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ahmedash95/gatewaySDK"
)

var PORT string
var SERVICE_NAME string
var API_GATEWAY string
var CONTAINER_URL string

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
		panic(fmt.Sprintf("Cant't regsiter comments service:%s", err.Error()))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World from service: %s <br> %s \n", SERVICE_NAME, r.URL.String())
	})

	http.HandleFunc("/all", allCommentsHandler)

	http.ListenAndServe(":"+PORT, nil)
}

func allCommentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All Comments goes here")
}
