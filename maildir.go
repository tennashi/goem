package goem

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/tennashi/goem/maildir"
)

type Maildir struct {
	Name string
}

func NewMaildir(path string) (*Maildir, error) {
	if !maildir.IsMaildir(path) {
		return nil, fmt.Errorf("%v is not maildir", path)
	}
	return &Maildir{
		Name: filepath.Base(path),
	}, nil

}

func Maildirs(rootPath string) ([]Maildir, error) {
	dirsInfo, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}
	var maildirs []Maildir
	for _, dirInfo := range dirsInfo {
		if !dirInfo.IsDir() {
			continue
		}
		path := filepath.Join(rootPath, dirInfo.Name())
		md, err := NewMaildir(path)
		if err != nil {
			continue
		}
		maildirs = append(maildirs, *md)
	}
	return maildirs, nil
}
