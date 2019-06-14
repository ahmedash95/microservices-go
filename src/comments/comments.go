package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	auth "github.com/ahmedash95/authSDK"
	"github.com/gorilla/mux"
)

func GetPostComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	var comments []Comment
	GetDB().Where("post_id = ?", postID).Find(&comments)

	jsonResponse(w, comments, 200)
}
func CreateComment(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var comment Comment
	err := decoder.Decode(&comment)
	if err != nil {
		jsonResponse(w, err, 400)
		return
	}

	validation_errors := comment.Validate()
	if validation_errors != nil {
		jsonResponse(w, validation_errors, 422)
		return
	}

	user := auth.GetUser(r)
	comment.UserID = user.ID

	GetDB().Create(&comment)

	jsonResponse(w, comment, 200)
}
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, _ := strconv.Atoi(vars["id"])

	var comment Comment

	GetDB().Find(&comment, commentID)

	GetDB().Delete(&comment)

	jsonResponse(w, nil, 200)
}
func ShowUserComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var comments []Comment
	GetDB().Where("user_id = ?", userID).Find(&comments)

	jsonResponse(w, comments, 200)
}
