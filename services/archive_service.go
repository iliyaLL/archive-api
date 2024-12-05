package services

import (
	"archive/zip"
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"

	"github.com/iliyaLL/archive-api/models"
)

type ArchiveService interface {
	GetArchiveInfo(file *multipart.FileHeader) (*models.ArchiveInfo, error)
	CreateArchive(files []*multipart.FileHeader) ([]byte, error)
}

type archiveService struct{}

func NewArchiveService() ArchiveService {
	return &archiveService{}
}

func (s *archiveService) GetArchiveInfo(file *multipart.FileHeader) (*models.ArchiveInfo, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return nil, err
	}

	archiveInfo := &models.ArchiveInfo{
		Filename:    file.Filename,
		ArchiveSize: float64(file.Size),
		TotalSize:   0,
		TotalFiles:  0,
		Files:       make([]models.FileInfo, 0),
	}

	for _, f := range zipReader.File {
		mimeType := mime.TypeByExtension(filepath.Ext(f.Name))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		fileInfo := models.FileInfo{
			FilePath: f.Name,
			Size:     float64(f.UncompressedSize64),
			Mimetype: mimeType,
		}

		archiveInfo.Files = append(archiveInfo.Files, fileInfo)
		archiveInfo.TotalSize += fileInfo.Size
		archiveInfo.TotalFiles++
	}

	return archiveInfo, nil
}

func (s *archiveService) CreateArchive(files []*multipart.FileHeader) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return nil, err
		}

		dst, err := zipWriter.Create(file.Filename)
		if err != nil {
			src.Close()
			return nil, err
		}

		if _, err := io.Copy(dst, src); err != nil {
			src.Close()
			return nil, err
		}

		src.Close()
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
