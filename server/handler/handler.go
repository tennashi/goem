package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/tennashi/goem"
)

type handler struct {
	rootPath string
}

func New(rootPath string) *handler {
	return &handler{
		rootPath: rootPath,
	}

}

func (h *handler) ListMaildir(w http.ResponseWriter, r *http.Request) {
	mds, err := goem.Maildirs(h.rootPath)
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

func (h *handler) ListMail(w http.ResponseWriter, r *http.Request) {
	dirName := chi.URLParam(r, "dirName")
	subDirName := r.URL.Query().Get("sub_dir")
	if subDirName == "" {
		subDirName = "cur"
	}
	path := filepath.Join(h.rootPath, dirName)
	ms, err := goem.Mails(path, subDirName)
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

func (h *handler) GetMail(w http.ResponseWriter, r *http.Request) {
	dirName := chi.URLParam(r, "dirName")
	key := chi.URLParam(r, "key")

	path := filepath.Join(h.rootPath, dirName)
	m, err := goem.GetMail(path, key)
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
