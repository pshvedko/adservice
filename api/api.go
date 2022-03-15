package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type Service interface {
	Read([]string, []string, int, int, []string) (interface{}, error)
	Create(float32, string, string, []string) (interface{}, error)
}

type Api struct {
	Service
}

func New(service Service) *Api {
	return &Api{
		Service: service,
	}
}

type AddQuery struct {
	Price       float32  `schema:"price,required"`
	Subject     string   `schema:"subject,required"`
	Description string   `schema:"description,required"`
	Photo       []string `schema:"photo"`
}

func (q *AddQuery) validate() error {
	//TODO implement me
	return nil
}

func (a *Api) Add(w http.ResponseWriter, r *http.Request) {
	var q AddQuery
	if err := a.parseForm(r); err != nil {
		a.writeError(w, http.StatusInternalServerError, err)
	} else if err := a.decodeForm(&q, r); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
	} else if v, err := a.Service.Create(q.Price, q.Subject, q.Description, q.Photo); err != nil {
		a.writeError(w, http.StatusNotAcceptable, err)
	} else {
		a.writeJson(w, http.StatusCreated, v)
	}
}

type ListQuery struct {
	Limit  int `schema:"limit"`
	Offset int `schema:"offset"`
	Sort   css `schema:"sort"`
}

func (q *ListQuery) validate() error {
	//TODO implement me
	return nil
}

func (a *Api) List(w http.ResponseWriter, r *http.Request) {
	var q ListQuery
	if err := a.parseForm(r); err != nil {
		a.writeError(w, http.StatusInternalServerError, err)
	} else if err := a.decodeForm(&q, r); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
	} else if v, err := a.Service.Read(nil, nil, q.Limit, q.Offset, q.Sort); err != nil {
		a.writeError(w, http.StatusNotFound, err)
	} else {
		a.writeJson(w, http.StatusOK, v)
	}
}

type GetQuery struct {
	Id    []string `schema:"id,required"`
	Field css      `schema:"field"`
}

func (q *GetQuery) validate() error {
	//TODO implement me
	return nil
}

func (a *Api) Get(w http.ResponseWriter, r *http.Request) {
	var q GetQuery
	if err := a.parseForm(r); err != nil {
		a.writeError(w, http.StatusInternalServerError, err)
	} else if err := a.decodeForm(&q, r); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
	} else if v, err := a.Service.Read(q.Id, q.Field, 0, 0, nil); err != nil {
		a.writeError(w, http.StatusNotFound, err)
	} else {
		a.writeJson(w, http.StatusOK, reflect.ValueOf(v).Index(0).Interface())
	}
}

func (a *Api) parseForm(r *http.Request) (err error) {
	switch strings.Split(r.Header.Get("Content-Type"), "; ")[0] {
	case "multipart/form-data":
		err = r.ParseMultipartForm(1 << 20)
		if err != nil {
			return
		}
	}
	err = r.ParseForm()
	if err != nil {
		return
	}
	for k, v := range mux.Vars(r) {
		r.Form[k] = []string{v}
	}
	return
}

var d = schema.NewDecoder()

type validator interface {
	validate() error
}

func (a *Api) decodeForm(v validator, r *http.Request) (err error) {
	switch strings.Split(r.Header.Get("Content-Type"), "; ")[0] {
	case "application/json":
		err = json.NewDecoder(r.Body).Decode(v)
	default:
		err = d.Decode(v, r.Form)
	}
	if err == nil {
		return v.validate()
	}
	return
}

func (a *Api) writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	if err != nil {
		_, _ = fmt.Fprintln(w, strings.ToLower(err.Error()))
	}
}

func (a *Api) writeJson(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if obj != nil {
		_ = json.NewEncoder(w).Encode(obj)
	}
}

type css []string // comma separated string list

func (s *css) UnmarshalText(b []byte) error {
	*s = append(*s, strings.Split(string(b), ",")...)
	return nil
}

func (s css) MarshalText() ([]byte, error) {
	return []byte(strings.Join(s, ",")), nil
}
