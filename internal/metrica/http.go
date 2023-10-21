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

//	Validate проверяет правильно ли указаны поля в структуре,
//
// ID не пустой, а значения MType - один из counter или gauge
// в случае если структура не валидная - возвращает ошибку с описанием
// проблемы: почему именно структура невалидная
func (m Metrica) Validate() error {
	if m.ID == "" {
		return errors.New("структура Metrica должна иметь не пустое поле ID")
	}
	if m.MType == "counter" && m.Delta == nil {
		return errors.New("структура Metrica c MType==counter должна иметь поле Delta")
	}
	if m.MType == "gauge" && m.Value == nil {
		return errors.New("структура Metrica c MType==gauge должна иметь поле Value")
	}
	if m.MType != CounterName && m.MType != GaugeName {
		return errors.New("структура Metrica должна иметь типа метрики MType - один из counter или gauge")
	}
	return nil
}

// Metricer используется, чтобы передать метрики по сети. Все метрики должны поддерживать этот интерфейс
type Metrer interface {
	Metrica(id string) Metrica
	ParseMetrica(Metrica)
}
