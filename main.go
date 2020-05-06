package main

import (
	"go2020/Proxy/Http"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("invalid params!")
		return
	}
	url := os.Args[1]
	pxy := &Http.HttpProxyServer{Host: url}
	pxy.Listen()
}
