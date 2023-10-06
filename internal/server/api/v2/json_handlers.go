// Этот пакет содержит хендлеры нового образца,
// где мы передаем значения при помощи json-запросов
// по маршрутам /update и /value
package apiv2

import (
	"fmt"
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

// Storager это интерфейс к хранилищу, которое использует именно этот API. Таким образом мы делаем хранилище зависимым от
// API,  а не наоборот
type Storager interface {
	Counter(name string) (metrica.Counter, bool)
	SetCounter(string, metrica.Counter)
	ListCounters() []string
	Gauge(name string) (metrica.Gauge, bool)
	SetGauge(string, metrica.Gauge)
	ListGauges() []string
}

// API это колленция http.HanlderFunc, которые обращаются к единому хранилищу store
type API struct {
	store Storager
}

// New создает новую
func New(store Storager) API {
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
func (api API) UpdateWithJSON(w http.ResponseWriter, r *http.Request) {
	//todo mix of two handlers + save+ response
	m := metrica.Metrica{}
	err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		log.Error().Str("location", "json update handler").Msg("Cant unmarshal data in body")
		http.Error(w, "^0^ не могу размаршалить тело сообщения", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// изменяем хранилище и значения в переменной m
	err = updateStorageAndValue(api.store, &m)
	if err != nil {
		httpErrorWithLogging(w, http.StatusBadRequest, "Ошибка обновления хранилища: %v", err)
		return
	}

	_, _, err = easyjson.MarshalToHTTPResponseWriter(&m, w)
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
func (api API) ValueWithJSON(w http.ResponseWriter, r *http.Request) {
	//todo mix of two handlers + save+ response
	m := metrica.Metrica{}
	err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		log.Error().Str("location", "json update handler").Msg("Cant unmarshal data in body")
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// пробуем обновить метрики. Там где мы проверяем err, нас не интересуют ошибки где пустые поля Value или
	err, found := addValueFromStorage(api.store, &m)
	if err != nil {
		httpErrorWithLogging(w, http.StatusBadRequest, "Ошибка метрики %+v: %v", m, err) // todo, а как бы сделать так, чтобы %v подсвечивался
		return
	} else if !found {
		httpErrorWithLogging(w, http.StatusNotFound, "В хранилище не найдена метрика %v", m.ID)
		return
	}

	_, _, err = easyjson.MarshalToHTTPResponseWriter(&m, w)
	if err != nil {
		log.Error().Str("location", "json get handler(on return)").Msgf("Cant marshal a return for %v %v", m.MType, m.ID)
		http.Error(w, "cant marshal result", http.StatusInternalServerError)
		return
	}
}

// httpErrorWithLogging отправляет сообщение об ошибке, параллельно дублируя в журнал. Работает быстрее, чем просто две функции отдельно.
// Во-первых, конкатенация строк происходит при помощи Spfrintf, а не сложением, а во вторых один раз на два вызова: и логгера, и http.Error()
//
// w - responseWriter вашего HTTP хендлера
// statusCode - код ответа сервера, напр. 200, 400, http.StatusNotFound, http.StatusOK
// format, params - типичные параметры, как в функции Printf
func httpErrorWithLogging(w http.ResponseWriter, statusCode int, format string, params ...interface{}) {
	s := fmt.Sprintf(format, params...)
	log.Error().Str("location", "json update handler").Msg(s)
	http.Error(w, s, statusCode)
	// TODO
	//
	// Возможно это пока единственный повод держать кастомный логгер, чтобы в нем была функция типа withHttpError(w)
}

// addValueFromStorage добавляет структуре request поле со значением из хранилища store,
// будь то поле delta или value. Изменяет исходную структуру request в угоду скорости работы.
// Если метрика в хранилище не найдена, то возвращает falsе.
//
// m будет валидным только если err!=nil, иначе мы не можем полагаться что поля value или delta
// можно прочитать и они будут соответствовать значения из хранилища. Например, если
// m - невалидная структура, то m.Delta и m.Value будут содержать nil, и если
// счетчик с именем m.ID не найден в хранилище store, то то же самое
//
// found==true Тогда и только тогда, когда err!=nil; found=false тоже возможен когда err!=nil
// err!=nil тогда, когда структура m неправильно оформлена
func addValueFromStorage(store Storager, m *metrica.Metrica) (e error, found bool) {

	// TODO
	//
	// в этомй функции мы вынуждены создавать указатели на int или float, поэтому чуть выгодней конечно возвращать структуру целиком
	// или раньше инициализировать ссылки на float или int
	//
	err := m.Validate()
	if err != nil {
		return fmt.Errorf("полученная структура оформлена неправильно: %v", err), false
	}

	switch m.MType {
	case metrica.CounterName:
		var c metrica.Counter
		c, found = store.Counter(m.ID)
		if !found {
			return nil, false
		}
		m.Delta = new(int64) // todo попадёт в кучу(((
		*m.Delta = int64(c)

	case metrica.GaugeName:
		var g metrica.Gauge
		g, found = store.Gauge(m.ID)
		if !found {
			return nil, false
		}
		m.Value = new(float64)
		*m.Value = float64(g)
	}

	// Кажется тут очень важная деталь архитектураная. Если бы у меня хранилище содержало значения типа int64 и float64,
	// То я мог бы тут добавить ссылки прям на хранилище, то есть прям на переменные хранящиеся в идеале в стеке
	// Но поскольку я занимаюсь конвертированием типов, я так сделать не могу
	//
	// Именно поэтому в этом коде мне нужно добалвять проверку на found, а могу бы, например nil привязывать,
	// если бы storage.Storager мог бы вернуть мне (*int)(nil)
	//
	// Однако ни мапы(memStore) ни БД мне такой возможности не предоставит, вопрос снят

	return nil, found
}

// updateStorageAndValue обновляет и значение m и хранилища store  в соответствии со значением
func updateStorageAndValue(store Storager, m *metrica.Metrica) error {
	err := m.Validate()
	if err != nil {
		return fmt.Errorf("полученная структура оформлена неправильно: %v", err)
	}

	err = m.ContainsUpdate()
	if err != nil {
		return fmt.Errorf("структура %+v не содержит обновлений: %v", m, err)
	}

	switch m.MType {
	case metrica.CounterName:
		oldVal, _ := store.Counter(m.ID)
		val := oldVal + metrica.Counter(*m.Delta)
		store.SetCounter(m.ID, val)

		*m.Delta = int64(val)
	case metrica.GaugeName:
		store.SetGauge(m.ID, metrica.Gauge(*m.Value))
		// Обновлять значение не требуется
	}
	return nil

}
