package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"github.com/pshvedko/adservice/api"
	"github.com/pshvedko/adservice/service"
	"github.com/pshvedko/adservice/storage"
)

func main() {
	m, err := storage.New(&mgo.DialInfo{
		Addrs:     []string{"mongo:27017"},
		Database:  "",
		Username:  "",
		Password:  "",
		PoolLimit: 32,
	})
	if err != nil {
		log.Fatal(err)
	}
	s := service.New(m)
	a := api.New(s)
	r := mux.NewRouter()
	v := r.PathPrefix("/api/v1").Subrouter()
	v.HandleFunc("/", a.List).Methods(http.MethodGet)
	v.HandleFunc("/", a.Add).Methods(http.MethodPost)
	v.HandleFunc("/{id:[a-z0-9]+}", a.Get).Methods(http.MethodGet)
	r.Use(api.LogMiddleware)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	h := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		<-c
		if err := h.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
		signal.Stop(c)
		close(c)
	}()
	if err := h.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-c
	m.Close()
}
