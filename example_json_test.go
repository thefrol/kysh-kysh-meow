package kyshkyshmeow_test

import (
	"fmt"
	"log"
	"math/rand"

	kyshkyshmeow "github.com/thefrol/kysh-kysh-meow"
)

// Метрики можно отправлять пачкой, для этого их надо упаковать и
// типы C и G, представляющие собой метрики типа: счетчик и величина.
// и потом воспользоваться функцией BatchUpdate
func Example_batch() {
	// величина, которую я хочу передать
	v := rand.Float64()

	// Упаковываю
	g := kyshkyshmeow.G{
		ID:    "mygauge",
		Value: v,
	}

	// и увеличим счетчик на 10
	с := kyshkyshmeow.C{
		ID:    "mycounter",
		Delta: 10,
	}

	// мой сервер находится по адресу
	addr := "http://localhost:8089"

	// отправляем метрики
	err := kyshkyshmeow.BatchUpdate(addr, g, с)
	if err != nil {
		log.Fatalf("не могу обновить: %v", err)
	}

	fmt.Println("ok")

	// Output:
	// ok
}
