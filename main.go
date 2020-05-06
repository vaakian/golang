package main

import (
	"go2020/Proxy/Http"
	"go2020/Proxy/Socks5"
	"log"
	url2 "net/url"
	"os"
	"time"
)

type Proxy interface {
	Listen()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("invalid params!")
		return
	}
	url := os.Args[1]
	var pxy Proxy
	parsed, err := url2.Parse(url)
	if err != nil {
		log.Fatal("ivalid url")
	}
	switch parsed.Scheme {
	case "socks5":
		pxy = Socks5.NewSocksProxy5Server(parsed.Host, time.Second*5)
	case "http":
		pxy = Http.NewHttpProxyServer(parsed.Host)
	default:
		log.Fatal("unsupported protocol")
	}
	pxy.Listen()
}
