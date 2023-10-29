package config

import (
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_configureAgent(t *testing.T) {

	tests := []struct {
		name        string
		defaults    Agent
		env         map[string]string
		commandLine string
		wantCfg     Agent
		panic       bool
	}{
		{
			name:        "без параметров строки",
			defaults:    Agent{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         nil,
			commandLine: "agent",
			wantCfg:     Agent{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
		},
		{
			name:        "Указать сервер через командную строку",
			defaults:    Agent{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         nil,
			commandLine: "agent -a localhost:8092",
			wantCfg:     Agent{Addr: "localhost:8092", ReportInterval: 2, PollingInterval: 1},
		},
		{
			name:        "Указать ключ через командную строку",
			defaults:    Agent{Key: "will be rewrited"},
			env:         nil,
			commandLine: "agent -k abcde",
			wantCfg:     Agent{Key: "abcde"},
		},
		{
			name:        "Указать ключ через переменные окружения",
			defaults:    Agent{Key: "will be rewrited"},
			env:         map[string]string{"KEY": "qwerty"},
			commandLine: "agent -k abcde",
			wantCfg:     Agent{Key: "qwerty"},
		},
		{
			name:        "Указать интервалы через командную строку",
			defaults:    Agent{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         nil,
			commandLine: "agent -a localhost:8092 -r 10 -p 11",
			wantCfg:     Agent{Addr: "localhost:8092", ReportInterval: 10, PollingInterval: 11},
		},
		{
			name:        "Указать адрес через переменную окружения",
			defaults:    Agent{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         map[string]string{"ADDRESS": "localhost:8088"},
			commandLine: "agent -a localhost:8092",
			wantCfg:     Agent{Addr: "localhost:8088", ReportInterval: 2, PollingInterval: 1},
		},
		{
			name:        "Указать все через переменную окружения",
			defaults:    Agent{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         map[string]string{"ADDRESS": "localhost:8088", "REPORT_INTERVAL": "3", "POLLING_INTERVAL": "4"},
			commandLine: "agent -a localhost:8092",
			wantCfg:     Agent{Addr: "localhost:8088", ReportInterval: 3, PollingInterval: 4},
		},
		{
			name:        "Отрицательное значение интервала вызывает панику в командной строке",
			defaults:    Agent{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         map[string]string{"ADDRESS": "localhost:8088"},
			commandLine: "agent -r -1",
			panic:       true,
		},
		{
			name:        "Отрицательное значение интервала вызывает панику в командной строке 2",
			defaults:    Agent{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         map[string]string{"ADDRESS": "localhost:8088"},
			commandLine: "agent -p -1",
			panic:       true,
		},

		{
			name:        "Отрицательное значение интервала вызывает панику в переменной окружения 3",
			defaults:    Agent{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         map[string]string{"ADDRESS": "localhost:8088", "REPORT_INTERVAL": "-1"},
			commandLine: "agent",
			panic:       true,
		},
		{
			name:        "Отрицательное значение интервала вызывает панику в переменной окружения 3",
			defaults:    Agent{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         map[string]string{"ADDRESS": "localhost:8088", "POLLING_INTERVAL": "-1"},
			commandLine: "agent",
			panic:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// очистим стандартный флагсет, переопределение флагов запрещено. Так что
			// если тут не очищать обработчики, то будет тест паниковать
			flag.CommandLine = flag.NewFlagSet("new", flag.PanicOnError)

			// очистим переменные окружения и записываем новые
			os.Unsetenv("ADDRESS")
			os.Unsetenv("REPORT_INTERVAL")
			os.Unsetenv("POLLING_INTERVAL")
			os.Unsetenv("KEY")

			if tt.env != nil {
				for k, v := range tt.env {
					os.Setenv(k, v)
				}
			}

			// подменяем командную спроку прям в пакете os
			os.Args = nil
			os.Args = append(os.Args, strings.Split(tt.commandLine, " ")...)

			// отловим панику
			defer func() {
				r := recover()
				if r != nil {
					assert.True(t, tt.panic, "Panicked but should not")
					return
				}
				assert.False(t, tt.panic, "We not panicked but should")
			}()

			cfg := Agent{}
			cfg.MustConfigure(tt.defaults)
			//проведем конфигурацию
			assert.True(t, reflect.DeepEqual(tt.wantCfg, cfg), "Итоговая конфигурация не совпадает с ожидаемой")
		})
	}
}
