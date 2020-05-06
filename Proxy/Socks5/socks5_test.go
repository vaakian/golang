package Socks5

import (
	"log"
	"net"
	"testing"
	"time"
)

func Test_byte(t *testing.T) {
	//b := []byte{0x00, 0xff}
	//fmt.Println(string(binary.BigEndian.Uint16(b[len(b)-2:])))
	//fmt.Println(binary.BigEndian.Uint16(b))
	sp := NewSocksProxy5Server(":1080", time.Second*5)
	sp.Listen()
}

func Test_join(t *testing.T) {

	log.Println(net.JoinHostPort("baidu.com", "445"))
}
