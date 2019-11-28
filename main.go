package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type UrlInfo struct {
	Url string `json: url`
}

func saveUrl(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		// setup headers
		c.Header("Allow", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "origin, content-type, accept")
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	} else if c.Request.Method == "POST" {
		var u UrlInfo
		c.BindJSON(&u)
		fmt.Println(u)
		c.JSON(http.StatusOK, gin.H{
			"url": u.Url,
		})
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, string([]byte("**hello world, Max!**")))
	})

	router.POST("/api/pingpong", saveUrl)
	router.Run(":" + port)
}
