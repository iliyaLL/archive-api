package services

import (
	"bytes"
	"io"
	"mime/multipart"
	"strconv"

	"gopkg.in/gomail.v2"
)

type MailService interface {
	SendFile(file *multipart.FileHeader, emails []string) error
}

type mailService struct {
	host     string
	port     string
	username string
	password string
}

func NewMailService(host, port, username, password string) MailService {
	return &mailService{host, port, username, password}
}

func (s *mailService) SendFile(file *multipart.FileHeader, emails []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.username)
	m.SetHeader("To", emails...)
	m.SetHeader("Subject", "File Attachment")
	m.SetHeader("text/plain", "Please find the attached file")

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return err
	}

	m.Attach(file.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := io.Copy(w, bytes.NewReader(buf.Bytes()))
		return err
	}))

	port, err := strconv.Atoi(s.port)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(s.host, port, s.username, s.password)
	return d.DialAndSend(m)
}
