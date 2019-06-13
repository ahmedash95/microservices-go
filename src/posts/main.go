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
	DB_Init()
	RegisterAuthService()
	RegisterWebServer()
}

func RegisterAuthService() {
	auth.Init(os.Getenv("AUTH_SERVICE_URI"))
}

func RegisterWebServer() {
	rtr := mux.NewRouter()

	rtr.Use(AuthMiddleware)

	rtr.HandleFunc("/", GetLastPosts).Methods("GET")
	rtr.HandleFunc("/create", PostsCreateHandler).Methods("POST")
	rtr.HandleFunc("/{id:[0-9]+}", PostsUpdateHandler).Methods("POST")
	rtr.HandleFunc("/{id:[0-9]+}", ShowPostHandler).Methods("GET")
	rtr.HandleFunc("/user/{id:[0-9]+}", ShowUserPosts).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":"+PORT, nil)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUser(r)
		if !user.IS_VALID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
