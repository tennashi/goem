package maildir

import (
	"errors"
	"io/ioutil"
	"net/mail"
	"os"
	"path/filepath"
)

// SubDir is the subdirectory name.
type SubDir uint8

const (
	// SubDirUnknown is unknown directory.
	SubDirUnknown SubDir = iota
	// SubDirCur is the cur sub directory.
	SubDirCur
	// SubDirNew is the new sub directory.
	SubDirNew
	// SubDirTmp is the tmp sub directory.
	SubDirTmp
)

// NewSubDir is create SubDir instance.
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

// String is ...
func (s SubDir) String() string {
	switch s {
	case SubDirCur:
		return "cur"
	case SubDirNew:
		return "new"
	case SubDirTmp:
		return "tmp"
	default:
		return "unknown"
	}
}

// Maildir is ...
type Maildir struct {
	Path string
}

// New is ...
func New(path string) (*Maildir, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &Maildir{
		Path: path,
	}, nil
}

// Mail is ...
type Mail struct {
	Key     Key
	Message *mail.Message
}

// Mails is ...
func (md Maildir) Mails(s SubDir) ([]Mail, error) {
	if s == SubDirUnknown {
		return nil, errors.New("unknown sub directory")
	}
	keys, err := md.Keys(s)
	if err != nil {
		return nil, err
	}
	SortKey(keys)

	ms := make([]Mail, len(keys))
	for i, k := range keys {
		err := func(i int) error {
			f, err := md.openMail(&k)
			if err != nil {
				return err
			}
			defer f.Close()

			m, err := mail.ReadMessage(f)
			if err != nil {
				return err
			}
			ms[i] = Mail{Key: k, Message: m}
			return nil
		}(i)
		if err != nil {
			return nil, err
		}
	}
	return ms, nil
}

// Keys is ...
func (md Maildir) Keys(s SubDir) ([]Key, error) {
	path := filepath.Join(md.Path, s.String())
	rawKeys, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	keys := make([]Key, len(rawKeys))
	for i, rawKey := range rawKeys {
		key, err := ParseKey(rawKey.Name())
		if err != nil {
			return nil, err
		}
		key.subDir = s
		keys[i] = key
	}
	return keys, nil
}

// Mail is ...
func (md Maildir) Mail(key Key) (*Mail, error) {
	f, err := md.openMail(&key)
	if err != nil {
		return nil, err
	}

	m, err := mail.ReadMessage(f)
	if err != nil {
		return nil, err
	}
	return &Mail{
		Key:     key,
		Message: m,
	}, nil
}

func (md Maildir) openMail(key *Key) (*os.File, error) {
	var p string
	switch key.subDir {
	case SubDirCur:
		p = filepath.Join(md.Path, "cur", key.String())
	case SubDirNew:
		p = filepath.Join(md.Path, "new", key.String())
	default:
		p = filepath.Join(md.Path, "cur", key.String())
	}

	f, err := os.Open(p)
	if err != nil {
		if err != os.ErrNotExist {
			return nil, err
		}
		p = filepath.Join(md.Path, "new", key.String())
		f, err = os.Open(p)
		if err != nil {
			return nil, err
		}
		if key.subDir == SubDirUnknown {
			key.subDir = SubDirNew
		}
		return f, err
	}

	if key.subDir == SubDirUnknown {
		key.subDir = SubDirCur
	}
	return f, nil
}

// IsMaildir is ...
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
