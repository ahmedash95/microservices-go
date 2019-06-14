package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dbs map[string]*gorm.DB

func DB_Init() {
	dbs = make(map[string]*gorm.DB)
	createConncetionFor("auth_service")
	createConncetionFor("posts_service")
	createConncetionFor("comments_service")
}

func createConncetionFor(dbName string) {
	query := getConnectionQuery(dbName)
	db, err := gorm.Open("postgres", query)
	if err != nil {
		panic(fmt.Sprintf("failed to connect %s database: %s", dbName, err.Error()))
	}
	dbs[dbName] = db
}

func getConnectionQuery(dbName string) string {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPass)
}

func GetDB(name string) *gorm.DB {
	return dbs[name]
}
