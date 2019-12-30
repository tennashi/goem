package goem

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/tennashi/goem/maildir"
)

// MaildirRoot is ...
type MaildirRoot struct {
	path string
}

// NewMaildirRoot is ...
func NewMaildirRoot(path string) *MaildirRoot {
	return &MaildirRoot{
		path: path,
	}
}

// Maildirs is ...
func (r *MaildirRoot) Maildirs() ([]Maildir, error) {
	dirsInfo, err := ioutil.ReadDir(r.path)
	if err != nil {
		return nil, err
	}
	var maildirs []Maildir
	for _, dirInfo := range dirsInfo {
		if !dirInfo.IsDir() {
			continue
		}
		path := filepath.Join(r.path, dirInfo.Name())
		md, err := NewMaildir(path)
		if err != nil {
			continue
		}
		maildirs = append(maildirs, *md)
	}
	return maildirs, nil
}

func (r *MaildirRoot) maildirPath(mdName string) string {
	return filepath.Join(r.path, mdName)
}

// GetMails is ...
func (r *MaildirRoot) GetMails(mdName, subDirName string) ([]Mail, error) {
	path := r.maildirPath(mdName)
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
		mail := NewMail(m)
		mails[i] = *mail
	}
	return mails, nil
}

// GetMail is ...
func (r *MaildirRoot) GetMail(mdName, key string) (*Mail, error) {
	path := r.maildirPath(mdName)
	if !maildir.IsMaildir(path) {
		return nil, fmt.Errorf("%v is not maildir", path)
	}
	md, err := maildir.New(path)
	if err != nil {
		return nil, err
	}
	k, err := maildir.ParseKey(key)
	if err != nil {
		return nil, err
	}
	ml, err := md.Mail(k)
	if err != nil {
		return nil, err
	}
	return NewMail(*ml), nil
}

// Maildir is ...
type Maildir struct {
	Name string
}

// NewMaildir is ...
func NewMaildir(path string) (*Maildir, error) {
	if !maildir.IsMaildir(path) {
		return nil, fmt.Errorf("%v is not maildir", path)
	}
	return &Maildir{
		Name: filepath.Base(path),
	}, nil
}
