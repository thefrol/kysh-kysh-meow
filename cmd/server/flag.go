package main

import (
	"flag"
	"fmt"
)

var addr *string = flag.String("a", ":8080", "[адрес:порт] устанавливает адрес сервера ")

func init() {

	flag.Usage = func() {
		print("server")
		flag.PrintDefaults()
		fmt.Println("^-^")
	}
}
