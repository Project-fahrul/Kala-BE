package service

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "PT. Makmur Subur Jaya <emailanda@gmail.com>"
const CONFIG_AUTH_EMAIL = "emailanda@gmail.com"
const CONFIG_AUTH_PASSWORD = "passwordemailanda"

type SmtpCredential struct {
	host         string
	port         int
	senderName   string
	authEmail    string
	authPassword string
}

var smtpCredential *SmtpCredential = nil

func SMTP_New() *SmtpCredential {
	if smtpCredential == nil {
		port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))

		if err != nil {
			port = 587
		}
		smtpCredential = &SmtpCredential{
			host:         os.Getenv("SMTP_HOST"),
			port:         port,
			senderName:   os.Getenv("SMTP_NAME"),
			authEmail:    os.Getenv("SMTP_AUTH_EMAIL"),
			authPassword: os.Getenv("SMTP_AUTH_PASSWORD"),
		}
	}
	return smtpCredential
}

func (s *SmtpCredential) SendConfirmMail(title string, to string, message string) error {
	dest := []string{to}
	messageBody := s.composeMessageBody(dest, title, message)
	return s.sending(dest, messageBody)
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{
		Name:    String,
		Address: "",
	}
	return strings.Trim(addr.String(), " <>")
}

func (s *SmtpCredential) composeMessageBody(dest []string, subject string, msg string) string {

	header := make(map[string]string)
	header["From"] = s.senderName
	header["To"] = strings.Join(dest, ",")
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + composeMessage(msg)
	return message
}

func composeMessage(msg string) string {
	return fmt.Sprintf("<html><body><p>%s</p></body></html>", msg)
}

func (s *SmtpCredential) sending(dest []string, msg string) error {
	auth := smtp.PlainAuth("", s.authEmail, s.authPassword, s.host)
	smptpAddress := fmt.Sprintf("%s:%d", s.host, s.port)

	err := smtp.SendMail(smptpAddress, auth, s.senderName, dest, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}
