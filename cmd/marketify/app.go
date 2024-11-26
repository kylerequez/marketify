package main

import (
	"github.com/kylerequez/marketify/src/servers"
)

func main() {
	if err := servers.Init(); err != nil {
		panic(err)
	}
}
