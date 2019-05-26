package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

type UserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u UserPayload
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}

	user := User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}

	if FindByEmail(user.Email).ID > 0 {
		response := map[string]string{
			"status":  "error",
			"message": "Email already exists in database",
		}

		body, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(422)
		w.Write(body)
		return
	}

	user.Password, _ = generateFromPassword(u.Password)
	GetDB().Create(&user)

	response := map[string]string{
		"status":  "success",
		"message": "User has been created succesfully",
	}

	body, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func FindByEmail(email string) User {
	var user User
	GetDB().Where("email = ?", email).First(&user)
	return user
}

func HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u LoginPayload
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}

	user := FindByEmail(u.Email)
	if user.Email != u.Email {
		response := map[string]string{
			"status":  "error",
			"message": "Invalid credentials",
		}
		body, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(422)
		w.Write(body)
		return
	}

	match, _ := comparePasswordAndHash(u.Password, user.Password)
	if !match {
		response := map[string]string{
			"status":  "error",
			"message": "Invalid password",
		}
		body, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(422)
		w.Write(body)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		response := map[string]string{
			"status":  "error",
			"message": err.Error(),
		}
		body, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(body)
		return
	}

	response := map[string]string{
		"status": "success",
		"token":  tokenString,
	}
	body, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func HandleGetLoggedInUser(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("X-token")
	if tokenString == "" {
		response := map[string]string{
			"status":  "error",
			"message": "Invalid token",
		}
		body, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(422)
		w.Write(body)
		return
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		response := map[string]string{
			"status":  "error",
			"message": "Invalid token",
		}
		body, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(422)
		w.Write(body)
		return
	}

	var user User
	var u UserResponse
	GetDB().First(&user, "id = ?", claims["user_id"])
	u.ID = user.ID
	u.Email = user.Email
	u.Name = user.Name

	response := map[string]interface{}{
		"status": "success",
		"user":   u,
	}
	body, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
