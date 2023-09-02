package main

import (
	"fmt"
	"strings"
)

// ParseURL рабивает строку типа https://my.server/update/<тип>/<имя>/<значение>
// на тип, имя и значение и возвращает в виде  URLParams
func ParseURL(url string) (p URLParams, e error) {
	p = URLParams(strings.Split(url, "/"))
	if !p.IsValid() {
		return nil, fmt.Errorf("cant parse url")
	}
	return
}

// URLParams содержит строки, на которые был разбит URL
// тип, имя метрики, переданное значение
type URLParams []string

const (
	TypePosition = iota + 2
	NamePosition
	ValuePosition
)

const urlSubpartsCount = 5

func (p URLParams) IsValid() bool {
	return len(p) == urlSubpartsCount
}

func (p URLParams) Type() string {
	return p[TypePosition]
}

func (p URLParams) Name() string {
	return p[NamePosition]
}

func (p URLParams) Value() string {
	return p[ValuePosition]
}
