package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// Hasher задает используемую функцию хеширования, создается каждый
// раз новая под каждый запрос
var Hasher func() hash.Hash = sha256.New

// Encoder будет переводить хеш в строку, по умолчанию исползуется base64
var Encoder func([]byte) string = base64.StdEncoding.EncodeToString

func Bytes(data []byte, key []byte) (string, error) {
	h := hmac.New(Hasher, key)
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	return Encoder(h.Sum(nil)), nil
}

// TODO
//
// В целом можно было бы организовать как пакет base64, например:
// есть какой-то класс довольно универсальный, типа Signer
// но в реальности мы пользуемся его конкретной реализацией, например
// sign.Sha256.WithKey("")
// и можно передавать в настройках, например этот sign.Sha256,
// а может его даже можно создать как экземпляр вместе с ключом
