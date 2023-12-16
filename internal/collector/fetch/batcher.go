package fetch

import "github.com/thefrol/kysh-kysh-meow/internal/metrica"

// Batcher это интерфейс, из которого можно доставать метрики пачками,
// в таких упавовках хранятся полученные, но ещё не превращаенные в
// metrica.Metrica метрики из runtime.MemStats
type Batcher interface {
	ToTransport() []metrica.Metrica
}
