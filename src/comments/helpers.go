package main

import (
	"encoding/json"
	"net/http"
)

type H map[string]string

func jsonResponse(w http.ResponseWriter, i interface{}, code int) {
	body, _ := json.Marshal(i)
	w.WriteHeader(code)
	w.Write(body)
}
