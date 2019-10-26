package main

import (
	"os"

	"github.com/tennashi/goem/cli/goemd"
)

/*
type Maildir struct {
	Name string `json:"name"`
}

type Mail struct {
	Key         string `json:"key"`
	Subject     string `json:"subject"`
	ContentType string `json:"content-type"`
}

func (m MaildirRoot) ListMail(dirName, subDirName string) ([]Mail, error) {
	path := filepath.Join(m.path, dirName)
	md, err := maildir.New(path)
	if err != nil {
		return nil, err
	}
	subDir := maildir.NewSubDir(subDirName)
	msgs, err := md.Messages(subDir)
	if err != nil {
		return nil, err
	}

	mails := make([]Mail, len(msgs))

	for i, msg := range msgs {
		mail := Mail{
			Key:     msg.Key.Raw,
			Subject: goem_mail.Header(msg.Message.Header).Get("Subject"),
		}
		mails[i] = mail
	}

	return mails, nil
}

type Message struct {
	Header     Mail
	Body       string
	Multiparts []string
}

func (m MaildirRoot) GetMail(dirName, subDirName, key string) (Message, error) {
	path := filepath.Join(m.path, dirName)
	md, err := maildir.New(path)
	if err != nil {
		return Message{}, err
	}
	subDir := maildir.NewSubDir(subDirName)

	msg, err := md.GetMessageWithRawKey(key, subDir)
	if err != nil {
		return Message{}, err
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return Message{}, err
	}

	var bodies []string
	var body string
	if strings.HasPrefix(mediaType, "multipart/") {
		r := multipart.NewReader(msg.Body, params["boundary"])
		for {
			p, err := r.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return Message{}, err
			}

			b64Decoder := base64.NewDecoder(base64.StdEncoding, p)

			body, err := ioutil.ReadAll(b64Decoder)
			if err != nil {
				return Message{}, err
			}
			bodies = append(bodies, string(body))
		}
		body = bodies[0]
	} else {
		b, err := ioutil.ReadAll(msg.Body)
		if err != nil {
			return Message{}, err
		}
		body = string(b)
	}

	return Message{
		Header: Mail{
			Key:         key,
			ContentType: msg.Header.Get("Content-Type"),
			Subject:     goem_mail.Header(msg.Header).Get("Subject"),
		},
		Body:       body,
		Multiparts: bodies,
	}, nil
}
*/

func main() {
	os.Exit(goemd.Run(os.Args))
}
