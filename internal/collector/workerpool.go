package collector

import (
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
	var batch []metrica.Metrica = make([]metrica.Metrica, 0, MaxBatch)
	defer sendBatch(batch, url) // это конечно нужно тестировать какой именно батч он отправит ахах
	// мне нужно все то, что осталось после закрытия канала входного

	for v := range inCh {
		batch = append(batch, v)
		if len(batch) >= MaxBatch {
			sendBatch(batch, url)
			batch = batch[:0]
		}
		// можно улучшить через default
	}

}
