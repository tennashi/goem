package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tennashi/goem"
)

// Handler is ...
type Handler struct {
	mdr *goem.MaildirRoot
}

// New is ...
func New(mdr *goem.MaildirRoot) *Handler {
	return &Handler{
		mdr: mdr,
	}
}

// ListMaildir is ...
func (h *Handler) ListMaildir(w http.ResponseWriter, r *http.Request) {
	mds, err := h.mdr.Maildirs()
	if err != nil {
		responseErr(w, err, http.StatusInternalServerError)
		return
	}

	type resp struct {
		Name string `json:"name"`
	}
	res := make([]resp, len(mds))
	for i, m := range mds {
		res[i] = resp{
			Name: m.Name,
		}
	}
	responseJSON(w, res, http.StatusOK)
}

// ListMail is ...
func (h *Handler) ListMail(w http.ResponseWriter, r *http.Request) {
	dirName := chi.URLParam(r, "dirName")
	subDirName := r.URL.Query().Get("sub_dir")
	if subDirName == "" {
		subDirName = "cur"
	}
	ms, err := h.mdr.GetMails(dirName, subDirName)
	if err != nil {
		responseErr(w, err, http.StatusBadRequest)
		return
	}

	type resp struct {
		Key     string              `json:"key"`
		Subject string              `json:"subject"`
		Headers map[string][]string `json:"headers"`
	}
	res := make([]resp, len(ms))
	for i, m := range ms {
		res[i] = resp{
			Key:     m.Key.Raw,
			Subject: m.Subject,
			Headers: m.Headers.DecodeAll(),
		}
	}
	responseJSON(w, res, http.StatusOK)
}

// GetMail is ...
func (h *Handler) GetMail(w http.ResponseWriter, r *http.Request) {
	dirName := chi.URLParam(r, "dirName")
	key := chi.URLParam(r, "key")

	m, err := h.mdr.GetMail(dirName, key)
	if err != nil {
		responseErr(w, err, http.StatusInternalServerError)
		return
	}
	type resp struct {
		Key     string              `json:"key"`
		Subject string              `json:"subject"`
		Body    string              `json:"body"`
		Headers map[string][]string `json:"headers"`
	}
	b, err := ioutil.ReadAll(m.Body)
	if err != nil {
		responseErr(w, err, http.StatusInternalServerError)
		return
	}

	res := resp{
		Key:     key,
		Subject: m.Subject,
		Body:    string(b),
		Headers: m.Headers.DecodeAll(),
	}
	responseJSON(w, res, http.StatusOK)
}

func responseErr(w http.ResponseWriter, err error, status int) {
	type retError struct {
		Error string
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	r := retError{Error: err.Error()}
	json.NewEncoder(w).Encode(r)

}

func responseJSON(w http.ResponseWriter, r interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(r)
}
