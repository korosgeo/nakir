package main

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"

	_ "github.com/lib/pq"
)

var db *sql.DB

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

func saveImageHandler(c *gin.Context) {
	file, err := c.FormFile("image")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No image received"})
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	path := "./images/" + filename

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"message": "Cannot save file to the file system"})
		return
	}

	db.Exec("INSERT INTO image (path) VALUES ($1)", path)
	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded", "filename": filename})
}

func getAllImagesHandler(c *gin.Context) {
	var images []string
	rows, err := db.Query("SELECT path FROM image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while retrieving images"})
	}
	defer rows.Close()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to get the images"})
		return
	}

	for rows.Next() {
		var image string
		if err := rows.Scan(&image); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Error while retrieving an image"})
		}
		images = append(images, image)
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", pingHandler)
	r.POST("/image", saveImageHandler)
	r.GET("/image/all", getAllImagesHandler)

	return r
}

func main() {
	db = setupDbConnection()
	defer db.Close()

	runMigrations(db)

	r := SetupRouter()
	r.Run(":8080")
}
