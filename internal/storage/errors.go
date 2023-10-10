package storage

import "fmt"

var (
	ErrorMetricNotFound = fmt.Errorf("метрика не найдена в хранилише")
)

var ErrorRestoreFileNotExist = fmt.Errorf("файла, чтобы выгрузить хранилише, не существует")

// да уж, по разнородности ошибок складывается впечатление, что пакет немного ХРЕНОВО спроектировн
// сверху ошибки только от storage