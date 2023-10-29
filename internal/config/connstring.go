package config

import "strings"

// ConnectionString создан для хранения строки соединения с БД. При попытке вывести такую структуру
// в консоль, часть строки с паролем покроется звездочкой.
// Пока что поддерживает только фомат "host= user= ..."
type ConnectionString struct {
	s string
}

func (c ConnectionString) String() string {
	var result []string
	for _, ss := range strings.Split(c.s, " ") {
		if strings.HasPrefix(ss, "password") {
			continue
		}
		result = append(result, ss)
	}
	return strings.Join(result, " ")
}

func (c ConnectionString) Get() string {
	return c.s
}

// Set воплощает интерфейс Value для пакета flag
func (c *ConnectionString) Set(s string) error {
	c.s = s
	return nil
}

// UnmarshalText воплощает интерфейс TextUnmarshaler
// пакета carlos0/env
func (c *ConnectionString) UnmarshalText(text []byte) error {
	return c.Set(string(text))
}
