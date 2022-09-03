package service

import (
	"fmt"
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

func (s *SmtpCredential) SendConfirmMail(to string, message string) error {
	dest := []string{to}
	messageBody := s.composeMessageBody(dest, "TEST", message)
	return s.sending(dest, messageBody)
}

func (s *SmtpCredential) composeMessageBody(dest []string, subject string, msg string) string {
	body := "From: " + s.senderName + "\n" +
		"To: " + strings.Join(dest, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		composeMessage(msg)

	return body
}

func composeMessage(msg string) string {
	return fmt.Sprintf("<a href='%s'>Click to confirm</a>", msg)
}

func (s *SmtpCredential) sending(dest []string, msg string) error {
	auth := smtp.PlainAuth("", s.authEmail, s.authPassword, s.host)
	smptpAddress := fmt.Sprintf("%s:%d", s.host, s.port)

	err := smtp.SendMail(smptpAddress, auth, s.authEmail, dest, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}
