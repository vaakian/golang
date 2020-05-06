package udp2tcp

import (
	"fmt"
	"log"
	"net"
	"testing"
)

func Test_Server(t *testing.T) {
	UDPListen(1081)
}

func Test_Client(t *testing.T) {
	TCPListen(1080)
}

func Test_DialHttp(t *testing.T) {

	ln, err := net.Dial("tcp", "hw.z-os.cn:3002")
	if err != nil {
		log.Println(err)
		return
	}
	ln.Write([]byte("GET Http://baidu.com/ HTTP/1.1\nHost: z-os.cn\nUser-Agent: curl/7.55.1\nAccept: */*\nProxy-Connection: Keep-Alive\n\n\n"))
	buf := make([]byte, 4096)
	size, err := ln.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf[0:size]))
}
