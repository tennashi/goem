package goem

import (
	"io"

	"github.com/tennashi/goem/mail"
	"github.com/tennashi/goem/maildir"
)

// Mail is ...
type Mail struct {
	Key     maildir.Key
	Subject string
	Headers mail.Header
	Body    io.Reader
}

// NewMail is ...
func NewMail(m maildir.Mail) *Mail {
	return &Mail{
		Key:     m.Key,
		Subject: mail.Header(m.Message.Header).Get("subject"),
		Headers: mail.Header(m.Message.Header),
		Body:    m.Message.Body,
	}
}
