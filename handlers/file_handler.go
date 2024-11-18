package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/iliyaLL/archive-api/services"
	"mime"
	"net/http"
	"path/filepath"
)

var allowedArchiveTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/xml": true,
	"image/jpeg":      true,
	"image/png":       true,
}

type FileHandler struct {
	archiveService services.ArchiveService
}

func NewFileHandler(as services.ArchiveService) *FileHandler {
	return &FileHandler{
		archiveService: as,
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
		print(mimetype)
		if !allowedArchiveTypes[file.Header.Get("Content-Type")] {
			c.JSON(http.StatusBadGateway, gin.H{"error": "invalid file mime type"})
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

func MailFile(c *gin.Context) {

}
