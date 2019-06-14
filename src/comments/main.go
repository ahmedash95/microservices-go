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
		panic(fmt.Sprintf("Cant't regsiter comments service:%s", err.Error()))
	}

	DB_Init()
	RegisterAuthService()
	registerWebServer()
}

func RegisterAuthService() {
	auth.Init(os.Getenv("AUTH_SERVICE_URI"))
}

func registerWebServer() {
	rtr := mux.NewRouter()

	rtr.Use(AuthMiddleware)
	rtr.Use(JsonResponseHeader)

	rtr.HandleFunc("/{id:[0-9]+}", GetPostComments).Methods("GET")
	rtr.HandleFunc("/create", CreateComment).Methods("POST")
	rtr.HandleFunc("/delete/{id:[0-9]+}", DeleteComment).Methods("POST")
	rtr.HandleFunc("/user/{id:[0-9]+}", ShowUserComments).Methods("GET")

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
		next.ServeHTTP(w, r)
	})
}

func JsonResponseHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
