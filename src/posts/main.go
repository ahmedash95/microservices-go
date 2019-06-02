package main

import (
	"fmt"
	"net/http"
	"os"

	auth "github.com/ahmedash95/authSDK"
	"github.com/ahmedash95/gatewaySDK"
	"github.com/gorilla/mux"
)

var PORT string
var SERVICE_NAME string
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
		panic(fmt.Sprintf("Cant't regsiter posts service:%s", err.Error()))
	}
	RegisterAuthService()
	RegisterWebServer()
}

func RegisterAuthService() {
	auth.Init(os.Getenv("AUTH_SERVICE_URI"))
}

func RegisterWebServer() {
	rtr := mux.NewRouter()

	rtr.Use(RegisterHeadersMidlleware)

	rtr.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World from service: %s\n", SERVICE_NAME)
	})

	rtr.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Posts home page")
	})

	rtr.HandleFunc("/post/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		postId := mux.Vars(r)["id"]
		fmt.Fprintf(w, "Display Post: %s", postId)
	})

	rtr.HandleFunc("/create", PostsCreateHandler)

	http.Handle("/", rtr)
	http.ListenAndServe(":"+PORT, nil)
}

func RegisterHeadersMidlleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Instance", r.Header.Get("X-Origin-Host"))
		next.ServeHTTP(w, r)
	})
}
