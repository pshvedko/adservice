package api

import (
	"io"
	"log"
	"net/http"
	"time"
)

type logHandle struct {
	h http.Handler
	w http.ResponseWriter
	c int
	n int
	u int
	b io.ReadCloser
}

func (l *logHandle) Read(p []byte) (n int, err error) {
	n, err = l.b.Read(p)
	l.u += n
	return
}

func (l *logHandle) Close() error {
	return l.b.Close()
}

func (l *logHandle) Header() http.Header {
	return l.w.Header()
}

func (l *logHandle) Write(b []byte) (n int, err error) {
	n, err = l.w.Write(b)
	l.n += n
	return
}

func (l *logHandle) WriteHeader(c int) {
	l.c = c
	l.w.WriteHeader(c)
}

func (l *logHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	l.w = w
	l.b = r.Body
	r.Body = l
	l.h.ServeHTTP(l, r)
	log.Println(r.Method, r.URL, r.Proto, l.c, l.u, l.n, time.Now().Sub(t))
}

func LogMiddleware(h http.Handler) http.Handler {
	return &logHandle{h: h, c: http.StatusOK}
}
