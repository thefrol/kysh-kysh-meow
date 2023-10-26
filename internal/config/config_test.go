package config

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
		defaults    ServerConfig
		env         map[string]string
		commandLine string
		wantCfg     ServerConfig
		panic       bool
	}{
		{
			name:        "без параметров строки",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
			env:         nil,
			commandLine: "serv",
			wantCfg:     ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "restore установить флагом",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         nil,
			commandLine: "serv -r",
			wantCfg:     ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "restore отменить флагом",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
			env:         nil,
			commandLine: "serv -r=false",
			wantCfg:     ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "restore со значением через переменную окружения",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -r",
			wantCfg:     ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "Указать ключ через командную строку",
			defaults:    ServerConfig{Key: newSecret("will be rewrited"), Restore: true},
			env:         nil,
			commandLine: "serv -k abcde",
			wantCfg:     ServerConfig{Key: newSecret("abcde"), Restore: true},
		},
		{
			name:        "Указать ключ через переменные окружения",
			defaults:    ServerConfig{Key: newSecret("will be rewrited"), Restore: true},
			env:         map[string]string{"KEY": "qwerty"},
			commandLine: "serv -k abcde",
			wantCfg:     ServerConfig{Key: newSecret("qwerty"), Restore: true},
		},
		{
			name:        "командной строкой указать файл куда писать",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -f 12342",
			wantCfg:     ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "12342"},
		},
		{
			name:        "переменной окружения указать файл куда писать",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true", "FILE_STORAGE_PATH": "123"},
			commandLine: "serv -f 1234",
			wantCfg:     ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "123"},
		},
		{
			name:        "командной строкой указать пустую строчку для файла",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: `serv -f `,
			wantCfg:     ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: ""},
		},
		{
			name:        "переменной окружения указать пустую строку для файла",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true", "FILE_STORAGE_PATH": ""},
			commandLine: "serv -r",
			wantCfg:     ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: ""},
		},
		{
			name:        "командной строкой указать интервал записи",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -i 299",
			wantCfg:     ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 299, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "переменной окружения указать интервал записи",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true", "STORE_INTERVAL": "123"},
			commandLine: "serv -i 22",
			wantCfg:     ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 123, Restore: true, FileStoragePath: "/tmp/file"},
		},

		{
			name:        "отрицательный интервал записи в командной строке вызывает панику",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -i -22",
			panic:       true,
		},
		{
			name:        "отрицательный интервал записи в переменной окружения вызывает панику",
			defaults:    ServerConfig{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
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
			actual := MustConfigure(tt.defaults)

			assert.True(t, reflect.DeepEqual(tt.wantCfg, actual), "Итоговая конфигурация не совпадает с ожидаемой")

		})
	}

}

func newSecret(s string) Secret {
	return Secret{value: []byte(s)}
}
