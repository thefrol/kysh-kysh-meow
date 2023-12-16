package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// подписывать будем алгоритмом SHA256,
// а кодировать байты при помощи base64
var (
	hasher  = sha256.New
	encoder = base64.StdEncoding.EncodeToString
)

func Bytes(data []byte, key []byte) (string, error) {
	h := hmac.New(hasher, key)
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	return encoder(h.Sum(nil)), nil
}

// TODO
//
// В целом можно было бы организовать как пакет base64, например:
// есть какой-то класс довольно универсальный, типа Signer
// но в реальности мы пользуемся его конкретной реализацией, например
// sign.Sha256.WithKey("")
// и можно передавать в настройках, например этот sign.Sha256,
// а может его даже можно создать как экземпляр вместе с ключом
