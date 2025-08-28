package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"

	_ "github.com/lib/pq"
)

func setupDbConnection() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres password=qwerty sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Db is connected.")
	}

	return db
}

func runMigrations(db *sql.DB) {
	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatal(err)
	}
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", pingHandler)
	// r.POST("/photo")

	return r
}

func main() {
	db := setupDbConnection()
	defer db.Close()

	runMigrations(db)

	r := SetupRouter()
	r.Run(":8080")
}
