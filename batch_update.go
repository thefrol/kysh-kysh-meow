package kyshkyshmeow

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mailru/easyjson"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

const batchUrl = "/updates"

// BatchUpdate передает метрики на сервер addr пачкой в одном запросе.
// Предварительно метрики надо упаковать в типы С и G
//
//	BatchUpdate("localhost:8080",
//		kyshkyshmeow.G{ID:"my_id",Value:20.1},
//		kyshkyshmeow.C{ID:"my_counter_id",Delta:20})
func BatchUpdate(addr string, metrics ...metrer) error {
	var b = make(metrica.Metricas, len(metrics))

	for i, m := range metrics {
		b[i] = m.toM()
	}

	data, err := easyjson.Marshal(&b)
	if err != nil {
		return fmt.Errorf("batch_update: %w", err)
	}

	url, err := url.JoinPath(addr, batchUrl)
	if err != nil {
		return fmt.Errorf("batch_update: %w", err)
	}

	body := bytes.NewBuffer(data)
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return fmt.Errorf("batch_update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorString, err := io.ReadAll(resp.Body)
		if err != nil {
			errorString = []byte("cant read body of response")
			// видимо вот для таких случаев надо осознанно работать с ошибками
		}

		return fmt.Errorf("batch_update: update_error: return code=%v reason=%v", resp.StatusCode, string(errorString))
	}

	return nil
}
