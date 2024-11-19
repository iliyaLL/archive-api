package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/iliyaLL/archive-api/handlers"
	"github.com/iliyaLL/archive-api/services"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func createTestServer() *gin.Engine {
	err := godotenv.Load("../.env")
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
	r.POST("/api/mail/file", fileHandler.SendFileEmail)
	r.POST("/api/archive/information", fileHandler.GetArchiveInfo)
	r.POST("/api/archive/files", fileHandler.CreateArchive)

	return r
}

func createMultipartRequest(uri string, fileFieldName string, filePath string, otherFields map[string]string) (*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileFieldName, file.Name())
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for key, val := range otherFields {
		_ = writer.WriteField(key, val)
	}

	writer.Close()
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestGetArchiveInformation(t *testing.T) {
	router := createTestServer()

	req, err := createMultipartRequest("/api/archive/information", "file", "./test_files/zip.zip", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, w.Code)
	}

	expected := "zip.zip"
	if !strings.Contains(w.Body.String(), expected) {
		t.Errorf("Expected response to contain %q, but got %s", expected, w.Body.String())
	}
}

func TestCreateArchive(t *testing.T) {
	router := createTestServer()

	req, err := createMultipartRequest("/api/archive/files", "files[]", "./test_files/nord.png", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/zip" {
		t.Errorf("Expected Content-Type application/zip, but got %s", contentType)
	}
}

func TestSendFileViaEmail(t *testing.T) {
	router := createTestServer()

	req, err := createMultipartRequest("/api/mail/file", "file", "./test_files/test.docx", map[string]string{
		"emails": "test1@test1.com,test2@test2.com",
	})
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, w.Code)
	}
}
