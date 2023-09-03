package main

import (
	"fmt"
	"runtime"
)

func main() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	ma, _ := getFieldsFloat(m)
	fmt.Printf("%+v", ma)
}
