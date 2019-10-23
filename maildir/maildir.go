package maildir

import (
	"errors"
	"io/ioutil"
	"net/mail"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Maildir map[SubDir]string

type SubDir uint8

const (
	SubDirCur SubDir = iota
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
		return SubDirCur
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

type FlagType uint8

const (
	_ FlagType = iota
	FlagTypeExperimental
	FlagTypeNormal
)

type Message struct {
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

func (md Maildir) Messages(s SubDir) ([]Message, error) {
	keys, err := md.GetKeys(s)
	if err != nil {
		return nil, err
	}
	SortKey(keys)

	ms := make([]Message, len(keys))
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
			msg := Message{Key: k, Message: *m}
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

func (md Maildir) GetMessage(key Key) (*mail.Message, error) {
	p := filepath.Join(md[key.s], key.Raw)
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	return mail.ReadMessage(f)
}

func (md Maildir) GetMessageWithRawKey(key string, s SubDir) (*mail.Message, error) {
	p := filepath.Join(md[s], key)
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	return mail.ReadMessage(f)
}

var (
	ErrCannotParse = errors.New("cannot parse")
)

type Key struct {
	s          SubDir
	Raw        string
	Second     uint
	DeliveryID ID
	HostName   string
	Params     map[string]string
	FlagType   FlagType
	Flags      []string
}

func ParseKey(str string) (Key, error) {
	k := Key{Raw: str}

	pieces := strings.SplitN(k.Raw, ".", 3)
	if len(pieces) < 3 {
		return Key{}, ErrCannotParse
	}
	sec, err := strconv.ParseUint(pieces[0], 10, 0)
	if err != nil {
		return Key{}, err
	}
	k.Second = uint(sec)

	k.DeliveryID, err = ParseID(pieces[1])

	additionals := strings.SplitN(pieces[2], ":", 2)

	h := strings.Split(additionals[0], ",")
	k.HostName = h[0]

	k.Params = make(map[string]string, len(h[1:]))
	for _, param := range h[1:] {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) < 2 {
			return Key{}, ErrCannotParse
		}
		k.Params[kv[0]] = kv[1]
	}

	flags := strings.Split(additionals[1], ",")
	ft, err := strconv.ParseUint(flags[0], 10, 8)
	if err != nil {
		return Key{}, err
	}
	k.FlagType = FlagType(ft)

	k.Flags = strings.Split(flags[1], "")

	return k, nil
}

func SortKey(ks []Key) {
	sort.Sort(keySlice(ks))
}

type keySlice []Key

func (ks keySlice) Len() int {
	return len(ks)
}

func (ks keySlice) Less(i, j int) bool {
	ksi := ks[i]
	ksj := ks[j]
	if ksi.Second > ksj.Second {
		return true
	}

	ids := idSlice{ksi.DeliveryID, ksj.DeliveryID}
	if ids.Less(0, 1) {
		return true
	}

	return false
}

func (ks keySlice) Swap(i, j int) {
	ks[i], ks[j] = ks[j], ks[i]
}

type ID struct {
	UNIXSeq     uint
	Boot        uint
	Urandom     uint
	Inode       uint
	Dev         uint
	MicroSecond uint
	PID         uint
	Seq         uint
}

func ParseID(str string) (ID, error) {
	checkNewFormat := func(r rune) bool {
		return !unicode.IsNumber(r)
	}

	if strings.Contains(str, "_") || strings.IndexFunc(str, checkNewFormat) == -1 {
		return parseOldFashionedID(str)
	}

	return parseID(str)
}

func parseID(str string) (ID, error) {
	id := ID{}

	for i := 0; i < len(str); {
		switch str[i] {
		case '#':
			id.UNIXSeq, i = parseValue(str, i)
		case 'X':
			id.Boot, i = parseValue(str, i)
		case 'R':
			id.Urandom, i = parseValue(str, i)
		case 'I':
			id.Inode, i = parseValue(str, i)
		case 'V':
			id.Dev, i = parseValue(str, i)
		case 'M':
			id.MicroSecond, i = parseValue(str, i)
		case 'P':
			id.PID, i = parseValue(str, i)
		case 'Q':
			id.Seq, i = parseValue(str, i)
		default:
			return ID{}, ErrCannotParse
		}
	}

	return id, nil
}

func parseValue(str string, i int) (uint, int) {
	idx := strings.IndexFunc(str[i+1:], func(r rune) bool {
		return !unicode.IsNumber(r)
	})
	if idx == -1 {
		idx = len(str) - i - 1
	}
	val, _ := strconv.ParseUint(str[i+1:i+idx+1], 10, 0)

	return uint(val), i + idx + 1
}

func parseOldFashionedID(str string) (ID, error) {
	id := ID{}

	ids := strings.SplitN(str, "_", 2)
	pid, err := strconv.ParseUint(ids[0], 10, 0)
	if err != nil {
		return ID{}, err
	}
	id.PID = uint(pid)

	if len(ids) < 2 {
		return id, nil
	}
	seq, err := strconv.ParseUint(ids[1], 10, 0)
	if err != nil {
		return ID{}, err
	}
	id.Seq = uint(seq)

	return id, nil
}

func SortID(ids []ID) {
	sort.Sort(idSlice(ids))
}

type idSlice []ID

func (ids idSlice) Len() int {
	return len(ids)
}

func (ids idSlice) Less(i, j int) bool {
	idi := ids[i]
	idj := ids[j]
	if idi.UNIXSeq > idj.UNIXSeq {
		return true
	}
	if idi.MicroSecond > idj.MicroSecond {
		return true
	}
	if idi.Seq > idj.Seq {
		return true
	}
	return false
}

func (ids idSlice) Swap(i, j int) {
	ids[i], ids[j] = ids[j], ids[i]
}
