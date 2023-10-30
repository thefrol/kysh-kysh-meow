package collector

import (
	"context"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

func FanIn(ctx context.Context, chs ...<-chan metrica.Metrica) chan metrica.Metrica {
	chMix := make(chan metrica.Metrica, len(chs)*generatorChannelSize)

	for _, ch := range chs {
		wg.Add(1)
		go func(ch <-chan metrica.Metrica) {
		loop:
			for {
				select {
				case v := <-ch:
					chMix <- v
				case <-ctx.Done():
					break loop
				}
				// надо придумывать как мне теперь о очереди всем каналы закрывать лол
			}
			wg.Done()
		}(ch)
	}

	// todo
	//
	// место под горутину остановки

	return chMix
}
