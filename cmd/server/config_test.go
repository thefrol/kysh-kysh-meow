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
		panic       bool
	}{
		{
			name:        "без параметров строки",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
			env:         nil,
			commandLine: "serv",
			wantCfg:     config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "restore установить флагом",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         nil,
			commandLine: "serv -r",
			wantCfg:     config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "restore отменить флагом",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
			env:         nil,
			commandLine: "serv -r=false",
			wantCfg:     config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "restore со значением через переменную окружения",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -r",
			wantCfg:     config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "Указать ключ через командную строку",
			defaults:    config{Key: "will be rewrited"},
			env:         nil,
			commandLine: "serv -k abcde",
			wantCfg:     config{Key: "abcde", Restore: true},
		},
		{
			name:        "Указать ключ через переменные окружения",
			defaults:    config{Key: "will be rewrited"},
			env:         map[string]string{"KEY": "qwerty"},
			commandLine: "serv -k abcde",
			wantCfg:     config{Key: "qwerty", Restore: true},
		},
		{
			name:        "командной строкой указать файл куда писать",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -f 12342",
			wantCfg:     config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "12342"},
		},
		{
			name:        "переменной окружения указать файл куда писать",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true", "FILE_STORAGE_PATH": "123"},
			commandLine: "serv -f 1234",
			wantCfg:     config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "123"},
		},
		{
			name:        "командной строкой указать пустую строчку для файла",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: `serv -f `,
			wantCfg:     config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: ""},
		},
		{
			name:        "переменной окружения указать пустую строку для файла",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true", "FILE_STORAGE_PATH": ""},
			commandLine: "serv -r",
			wantCfg:     config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: ""},
		},
		{
			name:        "командной строкой указать интервал записи",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -i 299", // todo может uint поставить???
			wantCfg:     config{Addr: "localhost:8081", StoreIntervalSeconds: 299, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "переменной окружения указать интервал записи",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true", "STORE_INTERVAL": "123"},
			commandLine: "serv -i 22",
			wantCfg:     config{Addr: "localhost:8081", StoreIntervalSeconds: 123, Restore: true, FileStoragePath: "/tmp/file"},
		},

		{
			name:        "отрицательный интервал записи в командной строке вызывает панику",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -i -22",
			panic:       true,
		},
		{
			name:        "отрицательный интервал записи в переменной окружения вызывает панику",
			defaults:    config{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true", "STORE_INTERVAL": "-2"},
			commandLine: "serv",
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
			os.Unsetenv("STORE_INTERVAL")
			os.Unsetenv("FILE_STORAGE_PATH")
			os.Unsetenv("RESTORE")
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

			//проведем конфигурацию
			actual := mustConfigure(tt.defaults)

			assert.True(t, reflect.DeepEqual(tt.wantCfg, actual), "Итоговая конфигурация не совпадает с ожидаемой")

		})
	}
}
