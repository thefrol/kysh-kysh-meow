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
	match := validPath.FindStringSubmatch(url)
	if match == nil {
		return nil, fmt.Errorf("url not matching its pattern %v", validPath)
	}

	u = &URLParams{}
	//#MENTOR как бы без мапы достать из regexp красиво
	for i, groupValue := range match {
		groupName := groupNames[i]
		switch groupName {
		case "type":
			u.Type = groupValue
		case "name":
			u.Name = groupValue
		case "value":
			u.Value = groupValue
		}
	}

	return
}
