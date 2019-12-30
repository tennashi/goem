package mail

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Header is ...
type Header mail.Header

// Get is ...
func (h Header) Get(key string) string {
	v := mail.Header(h).Get(key)
	return decodeHeader(v)
}

// Date is ...
func (h Header) Date() (time.Time, error) {
	return mail.Header(h).Date()
}

// AddressList is ...
func (h Header) AddressList(key string) ([]*mail.Address, error) {
	aList, err := mail.Header(h).AddressList(key)
	for _, a := range aList {
		a.Name = decodeHeader(a.Name)
	}
	return aList, err
}

// DecodeAll is ...
func (h Header) DecodeAll() Header {
	ret := make(map[string][]string, len(map[string][]string(h)))
	for key, values := range h {
		ret[key] = make([]string, len(values))
		for i, v := range values {
			ret[key][i] = decodeHeader(v)
		}
	}
	return ret
}

const (
	// ISO2022JPB is ...
	ISO2022JPB = "=?ISO-2022-JP?B?"
	// ISO2022JPQ is ...
	ISO2022JPQ = "=?ISO-2022-JP?Q?"
	// UTF8B is ...
	UTF8B = "=?UTF-8?B?"
	// UTF8Q is ...
	UTF8Q = "=?UTF-8?Q?"
	// SHIFTJISB is ...
	SHIFTJISB = "=?SHIFT_JIS?B?"
	// SHIFTJISQ is ...
	SHIFTJISQ = "=?SHIFT_JIS?Q?"
)

func decodeHeader(v string) string {
	fields := strings.Fields(v)
	out := make([]io.Reader, len(fields))
	for i, f := range fields {
		switch {
		case !strings.HasPrefix(f, "=?"):
			if i > 0 {
				f += " "
			}
			out[i] = strings.NewReader(f)
		case strings.HasPrefix(strings.ToUpper(f), ISO2022JPB):
			target := f[len(ISO2022JPB):strings.LastIndex(f, "?=")]
			r := strings.NewReader(target)

			b64Reader := base64.NewDecoder(base64.StdEncoding, r)
			iso2022jpDecoder := japanese.ISO2022JP.NewDecoder()
			out[i] = transform.NewReader(b64Reader, iso2022jpDecoder)
		case strings.HasPrefix(strings.ToUpper(f), ISO2022JPQ):
			target := f[len(ISO2022JPQ):strings.LastIndex(f, "?=")]
			r := strings.NewReader(target)

			qpReader := quotedprintable.NewReader(r)
			iso2022jpDecoder := japanese.ISO2022JP.NewDecoder()
			out[i] = transform.NewReader(qpReader, iso2022jpDecoder)
		case strings.HasPrefix(strings.ToUpper(f), SHIFTJISB):
			target := f[len(SHIFTJISB):strings.LastIndex(f, "?=")]
			r := strings.NewReader(target)

			b64Reader := base64.NewDecoder(base64.StdEncoding, r)
			shiftJISDecoder := japanese.ShiftJIS.NewDecoder()
			out[i] = transform.NewReader(b64Reader, shiftJISDecoder)
		case strings.HasPrefix(strings.ToUpper(f), SHIFTJISQ):
			target := f[len(SHIFTJISB):strings.LastIndex(f, "?=")]
			r := strings.NewReader(target)

			qpReader := quotedprintable.NewReader(r)
			shiftJISDecoder := japanese.ShiftJIS.NewDecoder()
			out[i] = transform.NewReader(qpReader, shiftJISDecoder)
		case strings.HasPrefix(strings.ToUpper(f), UTF8B):
			target := f[len(UTF8B):strings.LastIndex(f, "?=")]
			r := strings.NewReader(target)
			out[i] = base64.NewDecoder(base64.StdEncoding, r)
		case strings.HasPrefix(strings.ToUpper(f), UTF8Q):
			target := f[len(UTF8Q):strings.LastIndex(f, "?=")]
			r := strings.NewReader(target)
			out[i] = quotedprintable.NewReader(r)
		default:
			if i > 0 {
				f += " "
			}
			out[i] = strings.NewReader(f)
		}
	}

	b, _ := ioutil.ReadAll((io.MultiReader(out...)))
	return string(b)
}

// Message is ...
type Message struct {
	msg  *mail.Message
	mulR *multipart.Reader
}

// ReadMessage is ...
func ReadMessage(r io.Reader) (*Message, error) {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return nil, err
	}
	return &Message{msg: msg}, err
}

// SetBodyDecoder is ...
func (m *Message) SetBodyDecoder() {
	cType := m.msg.Header.Get("Content-Type")
	if cType == "" {
		return
	}
	mediaType, params, err := mime.ParseMediaType(cType)
	if err != nil {
		return
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		m.mulR = multipart.NewReader(m.msg.Body, params["boundary"])
	}
}
