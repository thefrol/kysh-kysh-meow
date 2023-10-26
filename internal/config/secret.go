package config

// Secret сделан специально для хранения деликатных данных, таких как пароли и ключ
// если случайно попадёт в Print(), то максимум в консоли появятся звездочки.
// Выдает значение только через функцию.
//
// Вообще изначально был план вообще сделать Secret как бы переименованием func() string
// но тогда надо постоянно проверять на nil, иначе все падает. Это не оч надежно канеш
type Secret struct {
	value []byte
}

func (k Secret) ValueFunc() func() []byte {
	return func() []byte { return k.value }
}

func (k Secret) IsEmpty() bool {
	return len(k.value) == 0
}

// String воплощает интерфейс fmt.Stringer,
// и возвращает звездочки, чтобы случайно в консоль не вывелось
func (k Secret) String() string {
	return "********"
}

// Set воплощает интерфейс Value для пакета flag
func (k *Secret) Set(s string) error {
	k.UnmarshalText([]byte(s))
	return nil
}

// UnmarshalText воплощает интерфейс TextUnmarshaler
// пакета carlos0/env
func (k *Secret) UnmarshalText(text []byte) error {
	k.value = text
	return nil
}
