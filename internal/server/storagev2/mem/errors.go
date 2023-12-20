package mem

import "errors"

var (
	ErrorNilStore = errors.New("обращение к хранилищу по nil указалелю")
	ErrorNilMap   = errors.New("мапа в хранилище не инищиализирована")
)
