// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge

package main

import (
	"net/http"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
	"github.com/thefrol/kysh-kysh-meow/internal/scheduler"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

var store storage.Storager

func init() {

}

var defaultConfig = config{
	Addr:                 ":8080",
	StoreIntervalSeconds: 300,
	FileStoragePath:      "/tmp/metrics-db.json",
	Restore:              true,
}

func main() {
	cfg := configure(defaultConfig)

	// создаем хранилище
	s := storage.New()
	store = wrapStorageWithWrite(time.Duration(cfg.StoreIntervalSeconds)*time.Second, s, func() {
		s.ToFile(cfg.FileStoragePath)
		ololog.Info().Msg("Хранилище сохранено в файл")
	})

	// подключаем сохранение хранилища на диск

	// Запускаем сервер в отдельной гоурутине
	ololog.Info().Msgf("^.^ Мяу, сервер запускается по адресу %v!", cfg.Addr)
	srv := http.Server{Addr: cfg.Addr, Handler: MeowRouter()}
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		ololog.Error().Msgf("^0^ не могу запустить сервер: %v \n", err)
	}

}

// CallBackStorage обертывает хранилище так, что при изменении значений вызывает специальную коллбек функцию
// SaveCallBack, которую можно назначить. Можно использовать для синхронной записи на диск, при изменениях значений
type CallBackStorage struct {
	storage.Storager
	SaveCallback func()
}

func (s CallBackStorage) SetCounter(name string, value metrica.Counter) {
	s.Storager.SetCounter(name, value)
	s.SaveCallback()
}
func (s CallBackStorage) SetGauge(name string, value metrica.Gauge) {
	s.Storager.SetGauge(name, value)
	s.SaveCallback()
}

// wrapStorageWithWrite заботится о том, чтобы данные из хранилища сохранялись на диск. Есть два режима:
// синхронная запись(writeInterval==0), и накапливаемая запись( writeInterval>0). И для записи
// в том и том случае будет использовать функция callback
func wrapStorageWithWrite(writeInterval time.Duration, s storage.Storager, writeCallback func()) storage.Storager {
	// Если нужна синхронная запись, значит оборачиваем хранилище в CallBackStorage.
	// И при каждом сохранении счетчика записываем все на диск
	if writeInterval == 0 {
		cbs := CallBackStorage{Storager: s}
		cbs.SaveCallback = writeCallback
		ololog.Info().Msg("Создано хранилище с синхронной записью на диск")
		return cbs
	}

	if writeInterval < 500*time.Millisecond {
		ololog.Warn().Str("location", "server storage wrapper").Msgf("Указана слишком быстрое время сохранения метрик %vс. Это может сказать на производительности", writeInterval.Seconds())
	}
	// в ином случае запускаем планировщик
	sc := scheduler.New()
	sc.AddJob(writeInterval, writeCallback)
	go sc.Serve(200 * time.Millisecond)
	return s
}
