package main

import (
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/jinzhu/gorm"
	"github.com/urfave/cli"
	pb "gopkg.in/cheggaaa/pb.v1"
)

/**
This file is responsible for creating users/posts/comments using http requests to the apis
*/

func main() {
	DB_Init()

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "fresh",
			Value: "false",
			Usage: "Purge all tables before seeding data",
		},
		cli.StringFlag{
			Name:  "count",
			Value: "10",
			Usage: "number of rows to seed",
		},
	}

	app := cli.NewApp()

	app.Flags = flags

	app.Commands = []cli.Command{
		{
			Name:   "all",
			Usage:  "seed all services posts/comments/users",
			Action: SeedAllAction,
			Flags:  flags,
		},
		{
			Name:   "posts",
			Usage:  "seed posts",
			Action: SeedPostsAction,
			Flags:  flags,
		},
		{
			Name:   "comments",
			Usage:  "seed comments",
			Action: SeedCommentsAction,
			Flags:  flags,
		},
		{
			Name:   "users",
			Usage:  "seed users",
			Action: SeedUsersAction,
			Flags:  flags,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func SeedAllAction(c *cli.Context) error {
	SeedUsersAction(c)
	SeedPostsAction(c)
	SeedCommentsAction(c)
	return nil
}
func SeedPostsAction(c *cli.Context) error {
	if c.String("fresh") == "true" {
		GetDB("posts_service").Unscoped().Delete(&Post{})
	}
	userIds := getUserIds()
	count, _ := strconv.Atoi(c.String("count"))
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		user := userIds[rand.Intn(len(userIds))]
		createPost(user.ID)
		bar.Increment()
	}
	bar.FinishPrint("")
	return nil
}
func SeedCommentsAction(c *cli.Context) error {
	if c.String("fresh") == "true" {
		GetDB("comments_service").Unscoped().Delete(&Comment{})
	}
	userIds := getUserIds()
	postIds := getPostIds()
	count, _ := strconv.Atoi(c.String("count"))
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		user := userIds[rand.Intn(len(userIds))]
		post := postIds[rand.Intn(len(postIds))]
		createComment(post.ID, user.ID)
		bar.Increment()
	}
	bar.FinishPrint("")
	return nil
}
func SeedUsersAction(c *cli.Context) error {
	if c.String("fresh") == "true" {
		GetDB("auth_service").Unscoped().Delete(&User{})
	}
	count, _ := strconv.Atoi(c.String("count"))
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		createUser()
		bar.Increment()
	}
	bar.FinishPrint("")
	return nil
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func createUser() {
	gofakeit.Seed(0)
	user := User{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: "123456",
	}
	GetDB("auth_service").Create(&user)
}

type Post struct {
	gorm.Model
	Title       string `json:"title"`
	Content     string `json:"content"`
	PublishDate string `json:"publish_date"`
	ImageURL    string `json:"image_url"`
	IsDraft     bool   `json:"is_draft"`
	UserID      int    `json:"user_id"`
}

func createPost(userID uint) {
	gofakeit.Seed(time.Now().UnixNano())
	post := Post{
		Title:       gofakeit.Sentence(4),
		Content:     gofakeit.Paragraph(10, 20, 100, " "),
		PublishDate: gofakeit.DateRange(time.Now().Add(-time.Hour*24*30), time.Now()).Format("2006-01-02 15:04"),
		ImageURL:    gofakeit.ImageURL(600, 300),
		IsDraft:     gofakeit.Bool(),
		UserID:      int(userID),
	}
	GetDB("posts_service").Create(&post)
}

type Comment struct {
	gorm.Model
	PostID  int    `json:"post_id"`
	Comment string `json:"comment"`
	UserID  int    `json:"user_id"`
}

func createComment(postID uint, userID uint) {
	gofakeit.Seed(time.Now().UnixNano())
	comment := Comment{
		PostID:  int(postID),
		UserID:  int(userID),
		Comment: gofakeit.Paragraph(3, 4, 5, ","),
	}
	GetDB("comments_service").Create(&comment)
}

func getUserIds() []User {
	var users []User
	GetDB("auth_service").Select("id").Find(&users)
	return users
}

func getPostIds() []Post {
	var posts []Post
	GetDB("posts_service").Select("id").Find(&posts)
	return posts
}
