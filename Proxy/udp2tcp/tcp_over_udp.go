package udp2tcp

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

// udp传输协议，http应用协议
// local:Http <--udp--> remote:Http

func UDPListen(Port int) {
	udpLn, udpErr := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: Port})
	if udpErr != nil {
		log.Fatal("unable to listen UDP: "+strconv.Itoa(Port), udpErr)
		return
	}
	udpBuf := make([]byte, 4096)
	for {
		// udpBuf存放客户端请求的数据
		size, remoteAddr, err := udpLn.ReadFromUDP(udpBuf)
		if err != nil {
			log.Fatal("fail to Accept connetion", err)
		} else {
			go HandleLocalRead(remoteAddr, udpBuf[0:size+1])
		}
	}
}
func HandleLocalRead(remoteAddr *net.UDPAddr, udpBuf []byte) {
	// dial local Http proxy
	log.Println("udpBuf: " + string(udpBuf))
	outConn, err := net.Dial("tcp", "hw.z-os.cn:3002")
	if err != nil {
		log.Println("unable to Dial local server")
		return
	}
	defer outConn.Close()
	log.Println("remote addr: " + remoteAddr.String())
	// 写客户端发送的请求到local Http proxy
	_, err = outConn.Write(udpBuf)
	if err != nil {
		log.Println("unable to send context to local server")
		return
	}
	// 当本地local接收到请求，将outConn的回复发送到remoteAddr

	for {
		readBuffer := make([]byte, 4096)
		// 读本地
		size, err := outConn.Read(readBuffer)
		if err != nil {
			log.Println("unable to read from httpServer", err)
			return
		}
		// 回写请求的数据，出口开放TCP就可以直接回写Tcp
		_, err = SendUDP(remoteAddr, readBuffer[0:size])
		fmt.Println(remoteAddr)
		if err != nil {
			log.Println("unable to sendUDP | ", err)
		}
	}
}
func SendUDP(addr *net.UDPAddr, ctx []byte) (*net.UDPConn, error) {
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: addr.IP, Port: addr.Port}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	conn.Write(ctx)
	log.Printf("<%s>\n", conn.RemoteAddr())
	return conn, nil
}

// 客户端

func TCPListen(Port int) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(Port))
	if err != nil {
		log.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			requestBuffer := make([]byte, 4096)
			size, err := conn.Read(requestBuffer)
			if err != nil {
				log.Println(err)
				return
			}
			// 远程udp地址，当收到udp包时，应该返回给tcp客户端
			remoteAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1081}
			udpConn, err := SendUDP(remoteAddr, requestBuffer[0:size])
			if err != nil {
				log.Println(err)
				return
			}
			defer udpConn.Close()
			replyBuffer := make([]byte, 4096)
			size, _, err = udpConn.ReadFromUDP(replyBuffer)
			log.Println("replayBuffer: ", string(replyBuffer))
			conn.Write(replyBuffer[0:size])
			log.Println("copy udp to tcp")
		}(conn)
	}
}
