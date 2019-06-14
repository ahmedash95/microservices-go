package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func DB_Init() {
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	query := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPass)
	fmt.Println(query)
	db, err = gorm.Open("postgres", query)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %s", err.Error()))
	}
}

func GetDB() *gorm.DB {
	return db
}
