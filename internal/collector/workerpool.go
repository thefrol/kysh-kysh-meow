package collector

import (
	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

const MaxBatch = 40

// worker считываем метрики и отправляет на сервер
//
// Подразумевается, что inCh это не какое-то постоянное соединение.
// То есть его создает тикер раз во сколько-то секунд
// func worker(inCh <-chan metrica.Metrica, url string) {
// 	batch := make([]metrica.Metrica, 0, MaxBatch)
// 	for v := range inCh {
// 		batch = append(batch, v)
// 	}
// 	sendBatch(batch, url)
// }
//
// это на случай создания и закрытия множества каналов под каждую отправку

func worker(inCh <-chan metrica.Metrica, url string, sema Semaphore) {
	batch := make([]metrica.Metrica, 0, MaxBatch)
	defer func() {
		sema.Acquire()
		sendBatch(batch, url)
		sema.Release()
		log.Info().Msg("Последний батч отправлен")
		// в данном случае мы не можем использовать неанонимную функцию
		// Ну можем... но тогда надо передавать ссылку на слайс, а не слайс
		//
		// если написать defer sendBatch(batch,url)
		// то go запомнит именно тот слайс, который мы передавали ранее
		// в момент создания(пустой), а нам нужен тот, что в момент завершения
	}()

	for v := range inCh {
		batch = append(batch, v)
		if len(batch) >= MaxBatch {
			sema.Acquire()
			sendBatch(batch, url)
			sema.Release()
			batch = batch[:0]
		}
		// можно улучшить через default

	}
}
