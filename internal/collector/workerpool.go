package collector

import (
	"sync"

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

// worker отправляет данные на сервер, ограничивая количество отправок семафором
// sema. Данные собираются из канала inCh и отправляются на сервер url, когда
// пачка становится размера MaxBatch.
//
// В идеале ещё бы отправлять по таймеру(todo), чтобы надолго не зависало, если
// данные во входном канале закончились.
func worker(inCh <-chan metrica.Metrica, url string, sema Semaphore, wg *sync.WaitGroup) {
	batch := make([]metrica.Metrica, 0, MaxBatch)
	defer wg.Done() // должен запуститься последним
	// важно сделать именно тут, чтобы это сработало уже после отправки последнего батча

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

// todo
//
// Для семантической красоты можно было бы сделать
// sema.Do(func (){
// 		sendBatch(batch)
// }
// код, который будет пропущен при срабатывании
//
// или сделать функцию типа semaSend(), или concurrentSend()
