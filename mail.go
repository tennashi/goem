package goem

import (
	"fmt"
	"io"

	"github.com/tennashi/goem/mail"
	"github.com/tennashi/goem/maildir"
)

type Mail struct {
	Key     maildir.Key
	Subject string
	Headers mail.Header
	Body    io.Reader
}

func NewMail(m maildir.Mail) *Mail {
	return &Mail{
		Key:     m.Key,
		Subject: mail.Header(m.Message.Header).Get("subject"),
		Headers: mail.Header(m.Message.Header),
		Body:    m.Message.Body,
	}
}

func Mails(path, subDirName string) ([]Mail, error) {
	if !maildir.IsMaildir(path) {
		return nil, fmt.Errorf("%v is not maildir", path)
	}
	md, err := maildir.New(path)
	if err != nil {
		return nil, err
	}
	sd := maildir.NewSubDir(subDirName)
	ms, err := md.Mails(sd)
	if err != nil {
		return nil, err
	}
	mails := make([]Mail, len(ms))
	for i, m := range ms {
		mails[i] = *NewMail(m)
	}
	return mails, nil
}

func GetMail(path, key string) (*Mail, error) {
	if !maildir.IsMaildir(path) {
		return nil, fmt.Errorf("%v is not maildir", path)
	}
	md, err := maildir.New(path)
	if err != nil {
		return nil, err
	}
	ml, err := md.GetMailWithRawKey(key)
	if err != nil {
		return nil, err
	}
	mail := NewMail(*ml)
	return mail, nil
}
