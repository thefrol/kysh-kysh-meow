// Этот пакет содержит хендлеры нового образца,
// где мы передаем значения при помощи json-запросов
// по маршрутам /update и /value
package apiv2

import (
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

var store storage.Storager

func SetStore(s storage.Storager) {
	store = s
}

// updateWithJSON обновляет значение счетчика JSON запросом. Читает из запроса тело в
// формате metrica.Metrica. Для счетчиков типа counter исползует поле delta и прибавляет к
// текущему значению, для счетчиков типа gauge заменяет текущее значение новым из поля Value.
// В ответ записывает структуру metrica.Metrica с обновленным значением
func UpdateWithJSON(w http.ResponseWriter, r *http.Request) {
	//todo mix of two handlers + save+ response
	m := metrica.Metrica{}
	err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		ololog.Error().Str("location", "json update handler").Msg("Cant unmarshal data in body")
		http.Error(w, "^0^ не могу размаршалить тело сообщения", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = m.Validate()
	if err != nil {
		ololog.Error().Str("location", "json update handler(validate metrica)").Msgf("Получена невалидная структура: %+v", err)
		http.Error(w, "Полученная структура оформлена неправильно: "+err.Error(), http.StatusBadRequest)
		return
	}

	switch m.MType {
	case metrica.CounterName:
		oldVal, _ := store.Counter(m.ID)
		val := oldVal + metrica.Counter(*m.Delta)
		store.SetCounter(m.ID, val) // todo: we need update counter!!!
	case metrica.GaugeName:
		store.SetGauge(m.ID, metrica.Gauge(*m.Value)) // todo точно надо в структере записать как gauge и counter поля
	default:
		ololog.Error().Str("location", "json update handler").Msgf("Cant update metric with type %v, no such metric", m.MType)
		http.Error(w, "unsupported type", http.StatusBadRequest)
		return
	}

	// returning можно сделать кой-то отдельной функицей
	var result metrica.Metrica
	var found bool
	switch m.MType {
	case metrica.CounterName:
		var c metrica.Counter
		c, found = store.Counter(m.ID) //todo меня бесит это возвращающее в два параметра, надо сделать функцию Lookup
		result = c.Metrica(m.ID)
	case metrica.GaugeName:
		var g metrica.Gauge
		g, found = store.Gauge(m.ID)
		result = g.Metrica(m.ID)
	default:
		ololog.Error().Str("location", "json update handler(on return)").Msgf("Cant update metric with type %v, no such metric", m.MType)
		http.Error(w, "unsupported type", http.StatusBadRequest)
		return
	}

	if !found {
		// todo такое ощущение, что нужен какой-то слой логики, везде одно и то же делаю как будто
		http.Error(w, "returned counter not found", http.StatusBadRequest)
	}

	_, _, err = easyjson.MarshalToHTTPResponseWriter(&result, w)
	if err != nil {
		ololog.Error().Str("location", "json update handler(on return)").Msgf("Cant marshal a return for %v %v", m.MType, m.ID)
		http.Error(w, "cant marshal result", http.StatusInternalServerError)
		return
	}
}

// value WithJSON возвращает значение счетчика JSON запросом. Читает из запроса тело в
// формате metrica.Metrica, у которого должны быть заполенены поля MType и ID,
// TДля счетчиков типа counter записывает значнием в поле delta,
// для счетчиков типа gauge в поле value. В ответ
// записывает структуру metrica.Metrica с обновленным значением
func ValueWithJSON(w http.ResponseWriter, r *http.Request) {
	//todo mix of two handlers + save+ response
	m := metrica.Metrica{}
	err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		ololog.Error().Str("location", "json update handler").Msg("Cant unmarshal data in body")
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var result metrica.Metrica
	var found bool
	switch m.MType {
	case metrica.CounterName:
		var c metrica.Counter
		c, found = store.Counter(m.ID) //todo меня бесит это возвращающее в два параметра, надо сделать функцию Lookup
		result = c.Metrica(m.ID)
	case metrica.GaugeName:
		var g metrica.Gauge
		g, found = store.Gauge(m.ID)
		result = g.Metrica(m.ID)
	default:
		ololog.Error().Str("location", "json value handler(on return)").Msgf("Cant get valye metric with type %v, no such metric", m.MType)
		http.Error(w, "unsupported type", http.StatusBadRequest)
		return
	}

	if !found {
		// todo такое ощущение, что нужен какой-то слой логики, везде одно и то же делаю как будто
		http.Error(w, "metrica not found", http.StatusNotFound)
		return
	}

	_, _, err = easyjson.MarshalToHTTPResponseWriter(&result, w)
	if err != nil {
		ololog.Error().Str("location", "json get handler(on return)").Msgf("Cant marshal a return for %v %v", m.MType, m.ID)
		http.Error(w, "cant marshal result", http.StatusInternalServerError)
		return
	}
}
