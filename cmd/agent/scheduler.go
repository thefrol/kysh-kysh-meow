package main

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

// Serve запускает запланированные работы, и занимает текущий поток
func (c *Chron) Serve(interval time.Duration) {
	if len(c.jobs) == 0 {
		fmt.Println("У планировщика нет задач. Ничего не выполняется...")
	}

	for {
		time.Sleep(interval)
		for _, job := range c.jobs {

			if job.Elapsed() {
				job.Run()
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

// Run Запускает функцию
func (j *Job) Run() {
	j.function()
	t := time.Now()
	j.lastCall = t
}
