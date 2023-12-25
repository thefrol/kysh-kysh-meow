package kyshkyshmeow_test

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/metricas"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router"
	"github.com/thefrol/kysh-kysh-meow/internal/server/storagev2/mem"
)

// инициализируем сервер, чтобы примеры срабатывали
// как тесты и всегда были актуальными
func init() {
	s := mem.MemStore{
		Counters: make(mem.IntMap, 50),
		Gauges:   make(mem.FloatMap, 50),
	}

	m := manager.Registry{
		Counters: &s,
		Gauges:   &s,
	}

	j := metricas.Manager{
		Registry: m,
	}

	api := router.API{
		Registry: m,
		Manager:  j,
	}

	h := api.MeowRouter()

	// запускаем сервер на две секунды
	go func() {
		s := http.Server{
			Addr:    ":8089",
			Handler: h,
		}
		go func() {
			err := s.ListenAndServe()
			if err != nil {
				log.Fatal(err)
			}
		}()

		// Дадим две секунды тестам завершиться
		time.Sleep(time.Second * 2)
		err := s.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()

}
