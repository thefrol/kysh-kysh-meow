package metrica

// Metrica это универсальный тип для передачи по HTTP
//
//easyjson:json
type Metrica struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge

	//todo
	//
	// Мне было бы приятно назвать эту структуру transport
}

//easyjson:json
type Metricas []Metrica

// Metricer используется, чтобы передать метрики по сети. Все метрики должны поддерживать этот интерфейс
type Metrer interface {
	Metrica(id string) Metrica
	ParseMetrica(Metrica)
}
