package gatewaySDK

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

var API_GATEWAY = os.Getenv("API_GATEWAY")

func SetGatewayURI(s string) {
	API_GATEWAY = s
}

type Service struct {
	NAME string `json:"service_name"`
	URL  string `json:"url"`
}

func RegisterService(s Service) (bool, error) {
	body, _ := json.Marshal(s)
	resp, err := http.Post(API_GATEWAY, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	return true, nil
}
