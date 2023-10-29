package config_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thefrol/kysh-kysh-meow/internal/config"
)

func Test_HideSecretFromConsole(t *testing.T) {
	original_key := "my_key"
	key := config.Secret{}
	key.Set(original_key)
	output_to_console := fmt.Sprint(key)
	assert.NotContains(t, output_to_console, original_key, "Вывод в консоль не должен содержать ключ исходный")
}

func Test_HidePasswordFromConsole(t *testing.T) {
	connstring := "host=localhost password=12342 dbname=123"
	cs := config.ConnectionString{}
	cs.Set(connstring)
	output_to_console := fmt.Sprint(cs)
	assert.NotContains(t, output_to_console, "password", "Вывод в консоль не должен содержать пароль")
}
