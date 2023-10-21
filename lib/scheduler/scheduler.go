package scheduler

import (
	"fmt"
	"time"
)

// Chron позволяет запускать функции с заданной переодичностью
// работы выполняются последовательно, стоит избегать очень долгих заданий
// Задания могут выполнится позднее, но никогда раньше
type Chron struct {
	jobs []*Job
}

// New создает новый планировщик Chron
func New() Chron {
	return Chron{}
}

// Serve запускает запланированные работы, и занимает текущий поток.
// Планировщик будет ходить по списку задач, с начала до конца,
// и потом опять. Если задачу уже можно запустить, то запускать.
// При этом делая паузу между такими циклами на interval
func (c *Chron) Serve(interval time.Duration) {
	if len(c.jobs) == 0 {
		fmt.Println("У планировщика нет задач. Ничего не выполняется...")
	}

	for {
		time.Sleep(interval)
		for _, job := range c.jobs {

			if job.Elapsed() {
				job.run()
			}
		}
	}
}

// AddJob добавляет задачу для выполенния интервально
func (c *Chron) AddJob(interval time.Duration, f func()) {
	now := time.Now()
	c.jobs = append(c.jobs, &Job{lastCall: now, interval: interval, function: f})
}

// Job хранит информацию о том, когда функцию нужно запускать и когда она была запущена последний раз
type Job struct {
	lastCall time.Time
	interval time.Duration
	function func()
}

// Elapsed возвращает true, если истек период ожидания, и функцию можно запускать заново
func (j *Job) Elapsed() bool {
	return time.Since(j.lastCall) > j.interval
}

// run Запускает функцию
func (j *Job) run() {
	t := time.Now()
	j.function()
	j.lastCall = t
}

// TODO
//
// Придумать способ остановки. Желательно, который занимает текущий поток. Чтобы мы точно знали, что у нас остановлено.
// Может быть после урока про контексты, я смогу припилить контекст к функции Server
