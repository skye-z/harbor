package main

import (
	"harbor/route"
)

func main() {
	route := route.NewRoute()
	route.Init()
	route.Run()
}
