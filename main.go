package main

import (
	"github.com/iliyaLL/archive-api/handlers"
	"github.com/iliyaLL/archive-api/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	fileHandler := handlers.NewFileHandler(
		services.NewArchiveService(),
		services.NewMailService(
			os.Getenv("SMTP_HOST"),
			os.Getenv("SMTP_PORT"),
			os.Getenv("SMTP_USERNAME"),
			os.Getenv("SMTP_PASSWORD"),
		),
	)

	r := gin.Default()
	r.HandleMethodNotAllowed = true

	archive := r.Group("/api/archive")
	{
		archive.POST("/information", fileHandler.GetArchiveInfo)
		archive.POST("/files", fileHandler.CreateArchive)
	}

	r.POST("api/mail/file", fileHandler.SendFileEmail)

	r.Run()
}
