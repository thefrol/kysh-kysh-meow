package manager

import (
	"fmt"
)

func IsValid(id string) error {
	if id == "" {
		return fmt.Errorf("пустой айдишник")
	}
	return nil
}
