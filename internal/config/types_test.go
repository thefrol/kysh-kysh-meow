package config_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thefrol/kysh-kysh-meow/internal/config"
)

func Test_HideSecretFromConsole(t *testing.T) {
	originalKey := "my_key"
	key := config.Secret{}
	key.Set(originalKey)
	consoleOutput := fmt.Sprint(key)
	assert.NotContains(t, consoleOutput, originalKey, "Вывод в консоль не должен содержать ключ исходный")
}

func Test_HidePasswordFromConsole(t *testing.T) {
	connString := "host=localhost password=12342 dbname=123"
	cs := config.ConnectionString{}
	cs.Set(connString)
	consoleOutput := fmt.Sprint(cs)
	assert.NotContains(t, consoleOutput, "password", "Вывод в консоль не должен содержать пароль")
}
