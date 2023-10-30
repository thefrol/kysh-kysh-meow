// Синглтон-пакет для эмуляции метрики RandomValue
package randomvalue

import (
	"math/rand"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

const (
	// название рамндомной метрики среди всех данных, что мы собираем
	IDRandomValue = "RandomValue"
)

// Get возвращает случайное число типа float64
func Get() metrica.Metrica {

	r := random.Float64()
	return metrica.Metrica{
		MType: metrica.GaugeName,
		ID:    IDRandomValue,
		Value: &r,
	}

}

var random *rand.Rand

func init() {
	s := rand.NewSource(int64(time.Now().Nanosecond()))
	random = rand.New(s)
}
