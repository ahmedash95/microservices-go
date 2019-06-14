package main

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	PostID  int    `json:"post_id"`
	Comment string `json:"comment"`
	UserID  int    `json:"user_id"`
}

func (p Comment) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.PostID, validation.Required, validation.Min(1)),
		validation.Field(&p.Comment, validation.Required, validation.Length(10, 100000)),
	)
}
