// Этот пакет содержит хендлеры нового образца,
// где мы передаем значения при помощи json-запросов
// по маршрутам /update и /value
package apiv2

import (
	"context"
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
)

// API это колленция http.HanlderFunc, которые обращаются к единому хранилищу store
type API struct {
	store api.Storager
}

// New создает новую
func New(store api.Storager) API {
	if store == nil {
		panic("Хранилище - пустой указатель")
	}
	return API{store: store}
}

// Не знаю, правильно ли это обновлять значение m в функциях AddValueFromStorage и UpdateStorageAndValue.
// Просто не хочется ещё раз обращаться к хранилищу, там проводить преобразование типов  и плодить указатели.
// Тут же у нас указатели в полях m.Delta и m.Value по заданию. Надо бы, конечно, провести эскейп анализ этого кода
//
// С точки зрения скорости такой код хорош, но кажется он не надежен и не прозрачен
//
// По иронии, с такими функциями код тепеь не дублируется, потому что при update не надо вызывать get из хранилиша.
// Но мне нравится, как они разгружают внешний вид самых хендлеров
//
// Есть ощущение, что можно объекдинить get и update в одну вспомогательную функцию. Это вроде как уберет дублирование кода, но
// лишь виртуально. Мы точно не знаем, добавятся не разойдутся ли эти функции в будущем. И вообще, мы нарушим принцип единственной ответственности.
// Кроме того, там добавится куча проверок; получим код, который я уверен, через неделю или две,  даже я не смогу понять, как работает.
// Эти проверки будут тормозить код. Так что с точки зрения и производительности и читаемости, мы все же разобьем на две функции, пусть они очень похожи,
// и изначальная задача была все же избавиться от дублирования.
//
// Нужно иметь в виду, что обработка этих хендлеров - особо чувствительная область. Представим, что мы пишем не в оперативную память,
// а в базу данных. Дополнительно вазывать чтение значение из базы - записали нормально ли? Не слишком ли это много времени займет, так что
// эти два хендлера работают в рачете на скорость, потому что тут это наиболее чувствительная область.

//

// updateWithJSON обновляет значение счетчика JSON запросом. Читает из запроса тело в
// формате metrica.Metrica. Для счетчиков типа counter исползует поле delta и прибавляет к
// текущему значению, для счетчиков типа gauge заменяет текущее значение новым из поля Value.
// В ответ записывает структуру metrica.Metrica с обновленным значением
func (i API) UpdateWithJSON(w http.ResponseWriter, r *http.Request) {
	//todo mix of two handlers + save+ response
	m := metrica.Metrica{}
	err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		log.Error().Str("location", "json update handler").Msg("Cant unmarshal data in body")
		http.Error(w, "^0^ не могу размаршалить тело сообщения", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// проверяем полученную структуру
	if err = m.Validate(); err != nil {
		api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Получена неправильно заполенная струкура %+v: %v", m, err)
		return
	}
	if err = m.ContainsUpdate(); err != nil {
		api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Получена неправильно заполенная струкура %+v: %v", m, err)
		return
	}

	// обновляем
	val, err := updateStorage(r.Context(), i.store, m)
	if err != nil {
		api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Ошибка обновления хранилища: %v", err)
		return
	}

	_, _, err = easyjson.MarshalToHTTPResponseWriter(&val, w)
	if err != nil {
		log.Error().Str("location", "json update handler(on return)").Msgf("Cant marshal a return for %v %v", m.MType, m.ID)
		http.Error(w, "cant marshal result", http.StatusInternalServerError)
		return
	}
}

// value WithJSON возвращает значение счетчика JSON запросом. Читает из запроса тело в
// формате metrica.Metrica, у которого должны быть заполенены поля MType и ID,
// TДля счетчиков типа counter записывает значнием в поле delta,
// для счетчиков типа gauge в поле value. В ответ
// записывает структуру metrica.Metrica с обновленным значением
func (i API) ValueWithJSON(w http.ResponseWriter, r *http.Request) {
	//todo mix of two handlers + save+ response
	m := metrica.Metrica{}
	err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		log.Error().Str("location", "json update handler").Msg("Cant unmarshal data in body")
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// проверяем полученную структуру
	if err = m.Validate(); err != nil {
		api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Получена неправильно заполенная струкура %+v: %v", m, err)
		return
	}

	val, err := GetStorage(r.Context(), i.store, m)
	if err != nil {
		if err == api.ErrorNotFoundMetric {
			api.HTTPErrorWithLogging(w, http.StatusNotFound, "В хранилище не найдена метрика %v", m.ID)
			return
		}
		api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Ошибка обновления метрики %+v: %v", m, err) // todo, а как бы сделать так, чтобы %v подсвечивался
		return
	}

	_, _, err = easyjson.MarshalToHTTPResponseWriter(&val, w)
	if err != nil {
		log.Error().Str("location", "json get handler(on return)").Msgf("Cant marshal a return for %v %v", m.MType, m.ID)
		http.Error(w, "cant marshal result", http.StatusInternalServerError)
		return
	}
}

func updateStorage(ctx context.Context, store api.Storager, upd metrica.Metrica) (newVal metrica.Metrica, err error) {
	switch upd.MType {
	case "counter":
		c, err := store.IncrementCounter(ctx, upd.ID, *upd.Delta)
		return metrica.Metrica{MType: upd.MType, ID: upd.ID, Delta: &c}, err // это получается отправится в хип
	case "gauge":
		g, err := store.UpdateGauge(ctx, upd.ID, *upd.Value)
		return metrica.Metrica{MType: upd.MType, ID: upd.ID, Value: &g}, err
	default:
		return
	}
}

func GetStorage(ctx context.Context, store api.Storager, req metrica.Metrica) (newVal metrica.Metrica, err error) {
	switch req.MType {
	case "counter":
		c, err := store.Counter(ctx, req.ID)
		return metrica.Metrica{MType: req.MType, ID: req.ID, Delta: &c}, err // это получается отправится в хип
	case "gauge":
		g, err := store.Gauge(ctx, req.ID)
		return metrica.Metrica{MType: req.MType, ID: req.ID, Value: &g}, err
	default:
		return
	}
}
