package main

import (
	"flag"
	"go2020/Proxy/Http"
	"go2020/Proxy/Socks5"
	"log"
	url2 "net/url"
	"time"
)

type Proxy interface {
	Listen()
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	var L = flag.String("L", "", "local listen host")
	var readBufferSize = flag.Uint("B", 128, "read buffer for socks5 server")
	flag.Parse()
	if *L == "" {
		//log.Println()
		log.Println("usage:")
		flag.PrintDefaults()
		log.Fatal("no schema specified")
	}

	var pxy Proxy
	parsed, err := url2.Parse(*L)
	if err != nil {
		log.Fatal("invalid url")
	}
	switch parsed.Scheme {
	case "socks5":
		if *readBufferSize < 64 || *readBufferSize > 65535 {
			log.Fatalf("invalid buffer size %d [64-65535]\n", *readBufferSize)
		}
		pxy = Socks5.NewSocksProxy5Server(parsed.Host, time.Second*5, uint16(*readBufferSize))

	case "http":
		pxy = Http.NewHttpProxyServer(parsed.Host)
	default:
		log.Fatal("unsupported protocol")
	}
	pxy.Listen()
}
