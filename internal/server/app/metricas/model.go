package metricas

import (
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

type Metrica = metrica.Metrica

// ValidateRequest валидирует входящую метрику
// у которой могут быть пустые поля Delta и Value
// возможно о таких случаях даже стоит логгировать как-то
func ValidateRequest(m Metrica) error {
	if m.ID == "" {
		return fmt.Errorf("%w: пустой айдишник", domain.ErrorBadID)
	}

	return nil
}

// todo по идее у нас еще должен быть тип
// данных типа RequestMetrica без всех этих ссылок

// todo
//
// в конце эта модель должна сюда переехать
