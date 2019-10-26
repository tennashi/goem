package maildir

import (
	"io/ioutil"
	"net/mail"
	"os"
	"path/filepath"
)

type SubDir uint8

const (
	SubDirUnknown SubDir = iota
	SubDirCur
	SubDirNew
	SubDirTmp
)

func NewSubDir(str string) SubDir {
	switch str {
	case "cur":
		return SubDirCur
	case "new":
		return SubDirNew
	case "tmp":
		return SubDirTmp
	default:
		return SubDirUnknown
	}
}

func (s SubDir) String() string {
	switch s {
	case SubDirCur:
		return "cur"
	case SubDirNew:
		return "new"
	case SubDirTmp:
		return "tmp"
	default:
		return ""
	}
}

type Maildir map[SubDir]string

type Mail struct {
	Key     Key
	Message mail.Message
}

func New(path string) (*Maildir, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &Maildir{
		SubDirCur: filepath.Join(path, "cur"),
		SubDirNew: filepath.Join(path, "new"),
		SubDirTmp: filepath.Join(path, "tmp"),
	}, nil
}

func (md Maildir) Mails(s SubDir) ([]Mail, error) {
	keys, err := md.GetKeys(s)
	if err != nil {
		return nil, err
	}
	SortKey(keys)

	ms := make([]Mail, len(keys))
	for i, k := range keys {
		err := func(i int) error {
			p := filepath.Join(md[s], k.Raw)
			f, err := os.Open(p)
			if err != nil {
				return err
			}
			defer f.Close()

			m, err := mail.ReadMessage(f)
			if err != nil {
				return err
			}
			msg := Mail{Key: k, Message: *m}
			ms[i] = msg
			return nil
		}(i)
		if err != nil {
			return nil, err
		}
	}
	return ms, nil
}

func (md Maildir) GetKeys(s SubDir) ([]Key, error) {
	rawKeys, err := ioutil.ReadDir(md[s])
	if err != nil {
		return nil, err
	}
	keys := make([]Key, len(rawKeys))
	for i, rawKey := range rawKeys {
		keys[i], err = ParseKey(rawKey.Name())
		keys[i].s = s
		if err != nil {
			return nil, err
		}
	}
	return keys, nil
}

func (md Maildir) GetMail(key Key) (*Mail, error) {
	p := filepath.Join(md[key.s], key.Raw)
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	m, err := mail.ReadMessage(f)
	if err != nil {
		return nil, err
	}
	return &Mail{
		Key:     key,
		Message: *m,
	}, nil
}

func (md Maildir) GetMailWithRawKey(key string) (*Mail, error) {
	k, err := ParseKey(key)
	if err != nil {
		return nil, err
	}
	return md.GetMail(k)
}

func IsMaildir(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	dirNames, err := f.Readdirnames(0)
	if err != nil {
		return false
	}
	var hasCur, hasNew, hasTmp bool
	for _, dirName := range dirNames {
		if dirName == "cur" {
			hasCur = true
		}
		if dirName == "new" {
			hasNew = true
		}
		if dirName == "tmp" {
			hasTmp = true
		}
	}
	return hasCur && hasNew && hasTmp
}
