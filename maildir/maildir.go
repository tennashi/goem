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

type Maildir struct {
	Path string
	keys struct {
		c []Key
		n []Key
	}
	SubDirs []Maildir
}

type FlagType uint8

const (
	_ FlagType = iota
	FlagTypeExperimental
	FlagTypeNormal
)

func Open(path string) (*Maildir, error) {
	m := &Maildir{
		Path: path,
	}

	curKeys, err := ioutil.ReadDir(filepath.Join(m.Path, "cur"))
	if err != nil {
		return nil, err
	}
	m.keys.c = make([]Key, len(curKeys))
	for i, key := range curKeys {
		m.keys.c[i], err = ParseKey(key.Name())
		if err != nil {
			return nil, err
		}
	}

	newKeys, err := ioutil.ReadDir(filepath.Join(m.Path, "new"))
	if err != nil {
		return nil, err
	}
	m.keys.n = make([]Key, len(newKeys))
	for i, key := range newKeys {
		m.keys.n[i], err = ParseKey(key.Name())
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (md Maildir) Messages() ([]*mail.Message, error) {
	ms := make([]*mail.Message, len(md.keys.c))
	SortKey(md.keys.c)
	for i, k := range md.keys.c {
		err := func(i int) error {
			p := filepath.Join(md.Path, "cur", k.Raw)
			f, err := os.Open(p)
			if err != nil {
				return err
			}
			defer f.Close()

			m, err := mail.ReadMessage(f)
			if err != nil {
				return err
			}
			ms[i] = m
			return nil
		}(i)
		if err != nil {
			return nil, err
		}
	}
	return ms, nil
}

func (md Maildir) NewMessages() ([]*mail.Message, error) {
	ms := make([]*mail.Message, len(md.keys.n))
	SortKey(md.keys.n)
	for i, k := range md.keys.n {
		f, err := os.Open(k.Raw)
		if err != nil {
			return nil, err
		}
		m, err := mail.ReadMessage(f)
		if err != nil {
			return nil, err
		}
		ms[i] = m
	}
	return ms, nil
}

var (
	ErrCannotParse = errors.New("cannot parse")
)

type Key struct {
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
