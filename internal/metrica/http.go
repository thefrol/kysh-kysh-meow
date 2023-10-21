package metrica

import (
	"errors"
)

// Metrica это универсальный тип для передачи по HTTP
type Metrica struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge

	//todo
	//
	// Мне было бы приятно назвать эту структуру transport
}

// Переменные выновятся в отдельный тип, чтобы не создавать их в рантайме постоянно и не очищать
var (
	ErrorNoID        = errors.New("структура Metrica c MType==gauge должна иметь поле Value")
	ErrorNoUpdates   = errors.New("структура Metrica c MType==counter должна иметь поле Delta с MType==gauge - поле Value")
	ErrorTypeInvalid = errors.New("структура Metrica должна иметь тип метрики MType - один из counter или gauge")
)

//	Validate проверяет правильно ли указаны поля в структуре,
//
// ID не пустой, а значения MType - один из counter или gauge
// в случае если структура не валидная - возвращает ошибку с описанием
// проблемы: почему именно структура невалидная
func (m Metrica) Validate() error {
	if m.ID == "" {
		return ErrorNoID
	}

	if m.MType != CounterName && m.MType != GaugeName {
		return ErrorTypeInvalid
	}
	return nil
}

// ContainsUpdate проверяет содержит ли структуда данные обновления. То есть value или delta должны быть не нулевыми ссылками
// Если обновления есть, то возврает err==nil
func (m Metrica) ContainsUpdate() error {
	if m.Value == nil && m.Delta == nil {
		return ErrorNoUpdates
	}
	return nil
}

// Metricer используется, чтобы передать метрики по сети. Все метрики должны поддерживать этот интерфейс
type Metrer interface {
	Metrica(id string) Metrica
	ParseMetrica(Metrica)
}
