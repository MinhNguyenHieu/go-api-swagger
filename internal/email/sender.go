package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type SMTPEmailSender struct {
	host        string
	port        int
	username    string
	password    string
	senderEmail string
}

func NewSMTPEmailSender(host string, port int, username, password, senderEmail string) *SMTPEmailSender {
	return &SMTPEmailSender{
		host:        host,
		port:        port,
		username:    username,
		password:    password,
		senderEmail: senderEmail,
	}
}

func (s *SMTPEmailSender) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := []byte("To: " + to + "\r\n" +
		"From: " + s.senderEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		body)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	err := smtp.SendMail(addr, auth, s.senderEmail, []string{to}, msg)
	if err != nil {
		log.Printf("Error sending email to %s: %v", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s with subject: %s", to, subject)
	return nil
}
