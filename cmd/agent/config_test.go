package main

import (
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_configure(t *testing.T) {

	tests := []struct {
		name        string
		defaults    config
		env         map[string]string
		commandLine string
		wantCfg     config
	}{
		{
			name:        "без параметров строки",
			defaults:    config{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         nil,
			commandLine: "agent",
			wantCfg:     config{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
		},
		{
			name:        "Указать сервер через командную строку",
			defaults:    config{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         nil,
			commandLine: "agent -a localhost:8092",
			wantCfg:     config{Addr: "localhost:8092", ReportInterval: 2, PollingInterval: 1},
		},
		{
			name:        "Указать интервалы через командную строку",
			defaults:    config{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         nil,
			commandLine: "agent -a localhost:8092 -r 10 -p 11",
			wantCfg:     config{Addr: "localhost:8092", ReportInterval: 10, PollingInterval: 11},
		},
		{
			name:        "Указать адрес через переменную окружения",
			defaults:    config{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         map[string]string{"ADDRESS": "localhost:8088"},
			commandLine: "agent -a localhost:8092",
			wantCfg:     config{Addr: "localhost:8088", ReportInterval: 2, PollingInterval: 1},
		},
		{
			name:        "Указать все через переменную окружения",
			defaults:    config{Addr: "localhost:8081", ReportInterval: 2, PollingInterval: 1},
			env:         map[string]string{"ADDRESS": "localhost:8088", "REPORT_INTERVAL": "3", "POLLING_INTERVAL": "4"},
			commandLine: "agent -a localhost:8092",
			wantCfg:     config{Addr: "localhost:8088", ReportInterval: 3, PollingInterval: 4},
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

			if tt.env != nil {
				for k, v := range tt.env {
					os.Setenv(k, v)
				}
			}

			// подменяем командную спроку прям в пакете os
			os.Args = nil
			os.Args = append(os.Args, strings.Split(tt.commandLine, " ")...)

			//проведем конфигурацию
			assert.True(t, reflect.DeepEqual(tt.wantCfg, configure(tt.defaults)), "Итоговая конфигурация не совпадает с ожидаемой")
		})
	}
}
