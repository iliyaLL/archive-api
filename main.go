package main

import (
	"github.com/iliyaLL/archive-api/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	r := gin.Default()
	r.HandleMethodNotAllowed = true

	archive := r.Group("/api/archive")
	{
		archive.POST("/information", handlers.ArchiveInformation)
		archive.POST("/files", handlers.ArchiveFiles)
	}

	r.POST("api/mail/file", handlers.MailFile)

	r.Run()
}
