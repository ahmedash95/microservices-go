package main

import (
	"encoding/json"
	"net/http"

	auth "github.com/ahmedash95/authSDK"
)

func PostsCreateHandler(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	j, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
