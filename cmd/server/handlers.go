package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

var store storage.MemStore

func init() {
	//Создать хранилище
	store = storage.New()
}

// updateCounter отвечает за маршрут, по которому будет обновляться счетчик типа counter
// иначе говоря за URL вида: /update/counter/<name>/<value>
// приходящее значение: int64
// поведение: складывать с предыдущим значением, если оно известно
func updateCounter(w http.ResponseWriter, r *http.Request, params URLParams) {
	value, err := strconv.ParseInt(params.Value, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "^0^ Ошибка значение, не могу пропарсить %v в %T", params.Value, value)
	}
	old, _ := store.Counter(params.Name)
	// по сути нам не надо обрабатывать случай, если значение небыло установлено. Оно ноль, прибавим новое значение и все четко
	new := old + storage.Counter(value)
	store.SetCounter(params.Name, new)

	fmt.Fprintf(w, "^.^ мур! Меняем Counter %v на %v. Новое значение %v", params.Name, value, new)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Printf("200(OK) at request to %v\n", r.URL.Path)
}

// updateGauge отвечает за маршрут, по которому будет обновляться метрика типа gauge
// иначе говоря за URL вида: /update/gauge/<name>/<value>
// приходящее значение: float64
// поведение: устанавливать новое значение
func updateGauge(w http.ResponseWriter, r *http.Request, params URLParams) {
	value, err := strconv.ParseFloat(params.Value, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "^0^ Ошибка значение, не могу пропарсить %v в %T", params.Value, value)
	}
	store.SetGauge(params.Name, storage.Gauge(value))

	fmt.Fprintf(w, "^.^ мур! меняем Gauge %v на %v.", params.Name, value)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Printf("200(OK) at request to %v\n", r.URL.Path)
}

// updateGauge отвечает за маршрут, по которому будет обновляться счетчик неизвестного типа
// без разбора возвращаем 400(Bad Request)
func updateUnknownType(w http.ResponseWriter, r *http.Request, params URLParams) {
	w.WriteHeader(http.StatusBadRequest) //обязательно за вызова w.Write() #INSIGHT добавить в README
	io.WriteString(w, "Фшшш! Я не знаю такой тип счетчика")
	w.Header().Add("Content-Type", "text/plain")
	fmt.Printf("400(BadRequest) at request to %v\n", r.URL.Path)
}
