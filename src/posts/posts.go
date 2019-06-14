package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	auth "github.com/ahmedash95/authSDK"
	"github.com/gorilla/mux"
)

func PostsCreateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var post Post
	err := decoder.Decode(&post)
	if err != nil {
		jsonResponse(w, err, 400)
		return
	}

	validation_errors := post.Validate()
	if validation_errors != nil {
		jsonResponse(w, validation_errors, 422)
		return
	}

	user := auth.GetUser(r)
	post.UserID = user.ID

	GetDB().Create(&post)

	jsonResponse(w, post, 200)
}

func PostsUpdateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newPost Post
	err := decoder.Decode(&newPost)
	if err != nil {
		jsonResponse(w, err, 400)
		return
	}

	validation_errors := newPost.Validate()
	if validation_errors != nil {
		jsonResponse(w, validation_errors, 422)
		return
	}

	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["id"])
	var post Post

	GetDB().Find(&post, postID)

	if post.ID != uint(postID) {
		jsonResponse(w, H{
			"message": "Invalid post id",
		}, 404)
		return
	}

	user := auth.GetUser(r)
	if post.UserID != user.ID {
		jsonResponse(w, H{
			"message": "This post doesn't belong to you",
		}, 404)
		return
	}

	post.Title = newPost.Title
	post.Content = newPost.Content
	post.PublishDate = newPost.PublishDate
	post.IsDraft = newPost.IsDraft
	post.ImageURL = newPost.ImageURL

	GetDB().Save(&post)

	jsonResponse(w, post, 200)
}

func ShowPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["id"])

	var post Post

	GetDB().Find(&post, postID)

	if post.ID != uint(postID) {
		jsonResponse(w, H{
			"message": "Post not found",
		}, 404)
		return
	}

	jsonResponse(w, post, 200)
}

func ShowUserPosts(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	var posts []Post

	var count int
	GetDB().Model(&Post{}).Where("user_id = ?", userID).Count(&count)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perpage := 2
	offset := 0
	isLastPage := 0
	if page < 2 {
		page = 0
	}
	offset = perpage * (page - 1)

	if offset+perpage >= count {
		isLastPage = 1
	}

	var query = GetDB()
	if user.ID == userID {
		query = query.Where("user_id = ?", userID)
	} else {
		query = query.Where("user_id = ? and is_draft = ?", userID, 1)
	}
	query.Order("publish_date desc").Offset(offset).Limit(perpage).Find(&posts)

	w.Header().Add("X-Posts-Total", strconv.Itoa(count))
	w.Header().Add("X-Last-Page", strconv.Itoa(isLastPage))

	jsonResponse(w, posts, 200)
}

func GetLastPosts(w http.ResponseWriter, r *http.Request) {
	var posts []Post

	var count int
	GetDB().Model(&Post{}).Count(&count)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perpage := 50
	offset := 0
	isLastPage := 0
	if page < 2 {
		page = 0
	}
	offset = perpage * (page - 1)

	if offset+perpage >= count {
		isLastPage = 1
	}

	GetDB().Order("publish_date desc").Where("is_draft = ?", 1).Offset(offset).Limit(perpage).Find(&posts)

	w.Header().Add("X-Posts-Total", strconv.Itoa(count))
	w.Header().Add("X-Last-Page", strconv.Itoa(isLastPage))

	jsonResponse(w, posts, 200)
}
