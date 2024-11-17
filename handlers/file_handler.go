package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/iliyaLL/archive-api/services"
	"net/http"
)

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
		c.JSON(http.StatusBadRequest, gin.H{"errro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

func ArchiveFiles(c *gin.Context) {

}

func MailFile(c *gin.Context) {

}
