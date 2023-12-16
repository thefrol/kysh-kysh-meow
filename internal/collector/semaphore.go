package collector

// Semaphore пропускает указанное количество горутин
// при помощи функции Acquire(), остальные ожадают,
// пока горутину не освободят место функцией Release()
//
// Используется для ограничения доступа к общим ресурсам,
// например, если мы хотим писать в базу данных не более
// двумя горутинами за раз
//
// sema:= NewSemaphore(2)
//
// sema.Acquire()
// db.Exec(...)
// sema.Release()
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore создает семафор, позволяющий
// проходит count горутинам через него
func NewSemaphore(count int) Semaphore {
	return Semaphore{
		ch: make(chan struct{}, count),
	}
}

// Acquire используется горутиной для прохода через семафор,
// если все пути заняты, то остановит горутину до освобождения
// семафора
func (sem Semaphore) Acquire() {
	sem.ch <- struct{}{}
}

// Release используется горутиной для выхода из семафора,
// освобождает место под горутину. Правило простое, если
// если воспользовался Acquire - не забудь и Release
func (sem Semaphore) Release() {
	<-sem.ch
}
