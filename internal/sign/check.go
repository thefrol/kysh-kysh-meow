package sign

import (
	"bytes"
	"crypto/hmac"
	"encoding/base64"
	"errors"
	"fmt"
)

var ErrorSignsNotEqual = errors.New("подписи не совпали")

// Check проверяет подпись
func Check(data []byte, key []byte, sign string) error {
	h := hmac.New(Hasher, key)
	_, err := h.Write(data)
	if err != nil {
		return fmt.Errorf("не могу записать хеш")
	}

	s := h.Sum(nil)

	decodedSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return fmt.Errorf("не могу декодировать подпись")
	}

	if !bytes.Equal(s, decodedSign) {
		return ErrorSignsNotEqual
	}
	return nil
}
