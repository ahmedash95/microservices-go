package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

/**
Instances is the container of the services and it's available instances
*/

type RegisterInstancePayload struct {
	ServiceName string `json:"service_name"`
	URL         string `json:"url"`
}

type Server struct {
	URL string
}

var Instances = make(map[string][]Server)

func main() {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "3001"
	}

	Instances[""] = append(Instances[""], Server{fmt.Sprintf("http://localhost:%s", PORT)})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url, _ := url.Parse(r.URL.String())

		paths := strings.Split(url.Path, "/")

		fmt.Println(paths[1])

		service := ""

		if len(paths) > 1 {
			service = paths[1]
		}

		fmt.Printf("Proxy service %s\n", service)

		serviceURL := Instances[""][0].URL

		if urls, Ok := Instances[service]; Ok {
			serviceURL = urls[0].URL
		}

		origin, _ := url.Parse(serviceURL)

		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}

		proxy := &httputil.ReverseProxy{Director: director}

		proxy.ServeHTTP(w, r)
	})

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
		panic(err)
	}

	serviceName := payload.ServiceName
	instanceURL := payload.URL

	newInstance := Server{instanceURL}

	Instances[serviceName] = append(Instances[serviceName], newInstance)
	fmt.Printf("Add new instance to service %s with url %s\n", serviceName, instanceURL)
}
