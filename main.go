package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)


type UrlInfo struct {
	Url string `json: url`
}


func saveUrl(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var u UrlInfo
		if c.Request.Method == "POST" {
			c.BindJSON(&u)
			fmt.Println(u)
		} else {
			return
		}

		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS urls (url text)"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
		}

		stmt, err := db.Prepare("INSERT INTO urls(url) VALUES(?)")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error preparing statement: %q", err)
			return
		}

		if _, err := stmt.Exec(u.Url); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error inserting url %s: %q", u.Url, err))
			return
		}
	}

	return gin.HandlerFunc(fn)
}

func showUrls(db *sql.DB) gin.HandlerFunc {
	fn := func (c *gin.Context) {
		rows, err := db.Query("SELECT * FROM urls")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading urls: %q", err))
			return
		}
		defer rows.Close()
		for rows.Next() {
			var url string
			if err := rows.Scan(&url); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning ticks: %q", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", url))
		}
	}
	return gin.HandlerFunc(fn)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, string([]byte("**hello world, Max!**")))
	})

	router.POST("/api/saveUrl", saveUrl(db))
	router.POST("/api/showUrls", showUrls(db))
	router.Run(":" + port)
}
