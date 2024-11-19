package handlers

import (
	"encoding/json"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iliyaLL/archive-api/services"
)

var allowedArchiveTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"text/xml; charset=utf-8": true,
	"image/jpeg":              true,
	"image/png":               true,
}

var allowedEmailTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/pdf": true,
}

type FileHandler struct {
	archiveService services.ArchiveService
	mailService    services.MailService
}

func NewFileHandler(as services.ArchiveService, ms services.MailService) *FileHandler {
	return &FileHandler{
		archiveService: as,
		mailService:    ms,
	}
}

func (h *FileHandler) GetArchiveInfo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	info, err := h.archiveService.GetArchiveInfo(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to format response"})
		return
	}

	c.Data(http.StatusOK, "application/json", response)
}

func (h *FileHandler) CreateArchive(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	files := form.File["files[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided"})
		return
	}

	for _, file := range files {
		mimetype := mime.TypeByExtension(filepath.Ext(file.Filename))
		if !allowedArchiveTypes[mimetype] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file mime type"})
			return
		}
	}

	archiveData, err := h.archiveService.CreateArchive(files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/zip", archiveData)
}

func (h *FileHandler) SendFileEmail(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	if !allowedEmailTypes[mime.TypeByExtension(filepath.Ext(file.Filename))] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file mime type"})
		return
	}

	emails := c.PostForm("emails")
	if emails == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no email address provided"})
		return
	}

	emailList := strings.Split(emails, ",")
	err = h.mailService.SendFile(file, emailList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
