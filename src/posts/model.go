package main

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Title       string `json:"title"`
	Content     string `json:"content"`
	PublishDate string `json:"publish_date"`
	ImageURL    string `json:"image_url"`
	IsDraft     bool   `json:"is_draft"`
	UserID      int    `json:"user_id"`
}

func (p Post) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, validation.Required, validation.Length(5, 50)),
		validation.Field(&p.Content, validation.Required, validation.Length(10, 100000)),
		validation.Field(&p.PublishDate, validation.NilOrNotEmpty, validation.Date("02-01-2006 15:04")),
		validation.Field(&p.ImageURL, validation.NilOrNotEmpty, is.URL),
	)
}
