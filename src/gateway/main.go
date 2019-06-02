package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type RegisterInstancePayload struct {
	ServiceName string `json:"service_name"`
	URL         string `json:"url"`
}

type Server struct {
	URL string
}

/** Instances is the container of the services and it's available instances*/
var Instances = make(map[string][]Server)

func main() {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "3001"
	}

	http.HandleFunc("/", ProxyHandler)

	http.HandleFunc("/instances", displayInstancesHandler)
	http.HandleFunc("/server-register", registerNewSerivce)

	http.ListenAndServe(":"+PORT, nil)
}

func displayInstancesHandler(w http.ResponseWriter, r *http.Request) {
	instances, _ := json.Marshal(Instances)
	w.Write(instances)
}

func registerNewSerivce(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload RegisterInstancePayload
	err := decoder.Decode(&payload)
	if err != nil {
		panic(fmt.Sprintf("Cant't regsiter the service:%s", err.Error()))
	}

	serviceName := payload.ServiceName
	instanceURL := payload.URL

	newInstance := Server{instanceURL}

	Instances[serviceName] = append(Instances[serviceName], newInstance)
	fmt.Printf("Add new instance to service %s with url %s\n", serviceName, instanceURL)
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	url, _ := url.Parse(r.URL.String())

	serviceURL, err := GetServiceURL(url)

	if err != nil {
		responseBody := map[string]string{
			"error": err.Error(),
		}
		body, _ := json.Marshal(responseBody)
		w.Write(body)
		return
	}

	origin, _ := url.Parse(serviceURL)

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}

	// remove service name from first param
	r.URL.Path = fmt.Sprintf("/%s", strings.Join(strings.Split(r.URL.Path, "/")[2:], "/"))

	proxy.ServeHTTP(w, r)
}

func GetServiceURL(u *url.URL) (string, error) {
	paths := strings.Split(u.Path, "/")

	if len(paths) < 1 {
		return "", errors.New("Invalid URL")
	}

	service := paths[1]

	urls, Ok := Instances[service]
	if !Ok {
		return "", errors.New("Service not found")
	}

	return urls[rand.Int()%len(urls)].URL, nil
}
