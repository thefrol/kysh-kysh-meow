package graceful

import (
	"context"
	"fmt"
	"os"
	"time"
)

// WithSignal возвращает контекст, который будет остановлен
// по запросу операционной системы, для таких сигналов как
// syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT
func WithSignal(ctx context.Context) context.Context {
	stoppable, cancel := context.WithCancel(ctx)
	sig := RequestStop()

	go func() {
		<-sig
		cancel()
		fmt.Println("горутина, созданная WithSignal(), остановлена")
	}()

	return stoppable
}

// SetForcedShutdown используется, чтобы завершить работу программы
// при истечении контекста ctx. Если в течение timeout секунд после
// истечения ctx, программа не завершилась, то это будет вызван
// sys.Exit(2)
func SetForcedShutdown(ctx context.Context, timeout time.Duration) {
	go func() {
		<-ctx.Done()
		terminationCtx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(timeout)) // интересно, если попробовать обеснуть уже истекший контекст?
		defer cancel()

		<-terminationCtx.Done()
		os.Exit(2)
	}()

	// todo
	//
	// Уже думаю, что семантически было бы прикольно
	// видеть shutdown.ForcedTimeout()
	// и shutdown.WithContext()
}
