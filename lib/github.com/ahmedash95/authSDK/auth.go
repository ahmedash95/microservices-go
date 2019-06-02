package authSDK

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var ENDPOINT_URI string

type User struct {
	ID       int    `json:"id"`
	NAME     string `json:"name"`
	EMAIL    string `json:"email"`
	IS_VALID bool
}

func Init(authServiceURI string) {
	ENDPOINT_URI = authServiceURI
}

func IsValidToken(token string) bool {
	fmt.Sprintf("%s with token: %s", ENDPOINT_URI+"/valid_token", token)
	client := &http.Client{}
	request, _ := http.NewRequest("POST", ENDPOINT_URI+"/valid_token", bytes.NewReader([]byte("")))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Token", token)
	resp, err := client.Do(request)

	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func ValidateRequestToken(r *http.Request) bool {
	token := r.Header.Get("X-token")
	return IsValidToken(token)
}

func GetUser(r *http.Request) User {
	invalidUser := User{IS_VALID: false, NAME: "a7aa"}
	valid := ValidateRequestToken(r)
	if !valid {
		fmt.Println("Invalid token")
		invalidUser.NAME = "invalid token"
		return invalidUser
	}

	client := &http.Client{}
	request, _ := http.NewRequest("POST", ENDPOINT_URI+"/user", bytes.NewReader([]byte("")))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Token", r.Header.Get("X-token"))
	resp, _err := client.Do(request)
	if _err != nil {
		fmt.Println("Invalid get user request")
		invalidUser.NAME = "invalid get user request"
		return invalidUser
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		fmt.Println(err)
		return invalidUser
	}
	u.IS_VALID = true

	return u
}
