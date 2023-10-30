package config

import (
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	a := ConnectionString{s: "124"}
	b := ConnectionString{s: "124"}
	assert.True(t, reflect.DeepEqual(a, b))
}

func Test_configure(t *testing.T) {

	tests := []struct {
		name        string
		defaults    Server
		env         map[string]string
		commandLine string
		wantCfg     Server
		wantErr     bool
	}{
		{
			name:        "без параметров строки",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
			env:         nil,
			commandLine: "serv",
			wantCfg:     Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "restore установить флагом",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         nil,
			commandLine: "serv -r",
			wantCfg:     Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "restore отменить флагом",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
			env:         nil,
			commandLine: "serv -r=false",
			wantCfg:     Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "restore со значением через переменную окружения",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -r",
			wantCfg:     Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "Указать ключ через командную строку",
			defaults:    Server{Key: newSecret("will be rewrited"), Restore: true},
			env:         nil,
			commandLine: "serv -k abcde",
			wantCfg:     Server{Key: newSecret("abcde"), Restore: true},
		},
		{
			name:        "Указать ключ через переменные окружения",
			defaults:    Server{Key: newSecret("will be rewrited"), Restore: true},
			env:         map[string]string{"KEY": "qwerty"},
			commandLine: "serv -k abcde",
			wantCfg:     Server{Key: newSecret("qwerty"), Restore: true},
		},
		{
			name:        "командной строкой указать файл куда писать",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -f 12342",
			wantCfg:     Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "12342"},
		},
		{
			name:        "переменной окружения указать файл куда писать",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true", "FILE_STORAGE_PATH": "123"},
			commandLine: "serv -f 1234",
			wantCfg:     Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: "123"},
		},
		{
			name:        "командной строкой указать пустую строчку для файла",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: `serv -f `,
			wantCfg:     Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: ""},
		},
		{
			name:        "переменной окружения указать пустую строку для файла",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true", "FILE_STORAGE_PATH": ""},
			commandLine: "serv -r",
			wantCfg:     Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: true, FileStoragePath: ""},
		},
		{
			name:        "командной строкой указать интервал записи",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -i 299",
			wantCfg:     Server{Addr: "localhost:8081", StoreIntervalSeconds: 299, Restore: true, FileStoragePath: "/tmp/file"},
		},
		{
			name:        "переменной окружения указать интервал записи",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true", "STORE_INTERVAL": "123"},
			commandLine: "serv -i 22",
			wantCfg:     Server{Addr: "localhost:8081", StoreIntervalSeconds: 123, Restore: true, FileStoragePath: "/tmp/file"},
		},

		{
			name:        "командной строкой указать строку подключения к бд",
			defaults:    Server{Restore: false, DatabaseDSN: newConnString("123")},
			env:         map[string]string{},
			commandLine: "serv -d qwqwqw",
			wantCfg:     Server{DatabaseDSN: newConnString("qwqwqw"), Restore: false},
		},
		{
			name:        "переменной указать строку подключения к бд",
			defaults:    Server{Restore: false},
			env:         map[string]string{"DATABASE_DSN": "1232"},
			commandLine: "serv -d qwqwqw",
			wantCfg:     Server{DatabaseDSN: newConnString("1232"), Restore: false},
		},

		{
			name:        "отрицательный интервал записи в командной строке вызывает панику",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true"},
			commandLine: "serv -i -22",
			wantErr:     true,
		},
		{
			name:        "отрицательный интервал записи в переменной окружения вызывает панику",
			defaults:    Server{Addr: "localhost:8081", StoreIntervalSeconds: 300, Restore: false, FileStoragePath: "/tmp/file"},
			env:         map[string]string{"RESTORE": "true", "STORE_INTERVAL": "-2"},
			commandLine: "serv",
			wantErr:     true,
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
				if r := recover(); r != nil {
					assert.Truef(t, tt.wantErr, "Panicked but should not:%v", r)
					return
				}
			}()

			//проведем конфигурацию
			actual := Server{}
			err := actual.Parse(tt.defaults)
			if tt.wantErr {
				assert.Error(t, err, "Должна быть ошибка, но ее нет")
				return
			}
			require.NoError(t, err)

			assert.True(t, reflect.DeepEqual(tt.wantCfg, actual), "Итоговая конфигурация не совпадает с ожидаемой")

		})
	}

}

func newSecret(s string) Secret {
	return Secret{value: []byte(s)}
}

func newConnString(s string) ConnectionString {
	return ConnectionString{s: s}
}
