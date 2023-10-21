// stats собирает информацию о использовании памяти, и сохраняет в хранилище
// так же имеет методы для управления счетчиком количества опросов PollCount
//
//	Fetch - основная функция пакета, достает из памяти нужные метрики и увеличивает
//
// счетчик PollCount
//
//	DropPoll - сбрасывает значение PollCount
//
// В целом этот пакет смотрит достаточно абстратно на метрики, просто сохраняет их в хранилище, дальше ему все равно
package stats

const (
	// название рамндомной метрики среди всех данных, что мы собираем
	randomValueName = "RandomValue"
	// Счеткич поличества опросов
	metricPollCount = "PollCount"
)
