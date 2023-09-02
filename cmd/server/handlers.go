package main

import (
	"fmt"
	"io"
	"net/http"
)

func updateCounter(w http.ResponseWriter, r *http.Request, params URLParams) {
	io.WriteString(w, "^.^ мур! Меняем Counter")
	w.Header().Add("Content-Type", "text/plain")
	fmt.Printf("200(OK) at request to %v\n", r.URL.Path)
}

func updateGauge(w http.ResponseWriter, r *http.Request, params URLParams) {
	io.WriteString(w, "^.^ мур! меняем Gauge")
	w.Header().Add("Content-Type", "text/plain")
	fmt.Printf("200(OK) at request to %v\n", r.URL.Path)
}

func updateUnknownType(w http.ResponseWriter, r *http.Request, params URLParams) {
	w.WriteHeader(http.StatusBadRequest) //обязательно за вызова w.Write() #INSIGHT добавить в README
	io.WriteString(w, "Фшшш! Я не знаю такой тип счетчика")
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("yoo", "yooo")
	fmt.Printf("400(BadRequest) at request to %v\n", r.URL.Path)
}
