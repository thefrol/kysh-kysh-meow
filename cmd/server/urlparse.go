package main

import (
	"fmt"
	"regexp"
)

// URLParams содержит параметры, на который был разбит URL
// тип, имя метрики, переданное значение
type URLParams struct {
	Type  string
	Name  string
	Value string
}

var validPath = regexp.MustCompile(`^/update/(?P<type>\w+)/(?P<name>\w+)/(?P<value>[^/]+)$`) //float=[+-]?[0-9]+(\.[0-9]+)?([Ee][+-]?[0-9]+)?

func ParseUrl(url string) (u *URLParams, e error) {
	groupNames := validPath.SubexpNames()
	match := validPath.FindAllString(url, -1)
	if match == nil {
		return nil, fmt.Errorf("url not matching its pattern %v", validPath)
	}

	u = &URLParams{}
	parsedUrl := make(map[string]string)
	for i, groupValue := range match {
		groupName := groupNames[i]
		parsedUrl[groupName] = groupValue
	}
	u.Type = parsedUrl["type"]
	u.Name = parsedUrl["name"]
	u.Value = parsedUrl["value"]
	return
}
