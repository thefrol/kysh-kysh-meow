package kyshkyshmeow_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// Самый простой спрособ отправить и получить метрики при помощи
// url запросов.
func Example_raw() {
	// счетчики, которые мы будем отправлять
	var (
		MyGauge     = 3.14
		MyGaugeName = "my_gauge"

		MyCounter     = 3
		MyCounterName = "my_counter"
	)

	// адрес сервера
	var (
		Addr = "localhost:8089"
	)

	// Чтобы отправить метрики, можно воспользоваться таким маршрутом
	// POST http://<my_server>/update/<тип метрики: gauge|counter>/<имя метрики>/<значение>

	// создаем запрос на сохранение метрики gauge
	url := fmt.Sprintf("http://%v/update/gauge/%v/%v", Addr, MyGaugeName, MyGauge)
	r, err := http.Post(url, "text/plain", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	// создаем запрос на сохранение счетчика
	url = fmt.Sprintf("http://%v/update/counter/%v/%v", Addr, MyCounterName, MyCounter)
	r, err = http.Post(url, "text/plain", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	// Получать метрики можно вот по такому запросу
	// GET http://<my_server>/update/<тип метрики: gauge|counter>/<имя метрики>

	// получаем величину
	url = fmt.Sprintf("http://%v/value/gauge/%v", Addr, MyGaugeName)
	r, err = http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	fmt.Println("my_gauge:", string(body))

	// получаем счетчик
	url = fmt.Sprintf("http://%v/value/counter/%v", Addr, MyCounterName)
	r, err = http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	body, err = io.ReadAll(r.Body)
	fmt.Println("my_counter:", string(body))

	// Output:
	// my_gauge: 3.14
	// my_counter: 3
}
