package Socks5

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type SocksProxy5Server struct {
	Host       string
	dialer     net.Dialer
	BufferSize uint16
}

func NewSocksProxy5Server(host string, timeout time.Duration, bufferSize uint16) *SocksProxy5Server {
	return &SocksProxy5Server{Host: host, dialer: net.Dialer{Timeout: timeout}, BufferSize: bufferSize}
}

func (sp *SocksProxy5Server) Listen() {
	// 服务器监听并把客户端请求交由handleClientRequest处理
	srvLn, err := net.Listen("tcp", sp.Host)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("socks5 proxy server listen on: [%s]\n", sp.Host)
	for {
		client, err := srvLn.Accept()
		if err != nil {
			log.Printf("proxySrv got an err while accepting [%s]\n", err)
			continue
		}
		go sp.handleClientRequest(client)
	}
}
func (sp *SocksProxy5Server) handleClientRequest(client net.Conn) {
	clientBuffer := make([]byte, sp.BufferSize)
	size, err := client.Read(clientBuffer)
	if err != nil {
		log.Printf("unable to read schema [%s]\n ", err)
		return
	}
	//clientBuffer[:size], []byte{0x05, bytesOfMethod, method}
	if clientBuffer[0] == 0x05 {
		client.Write([]byte{0x05, 0x00})
	} else {
		log.Printf(" unsupported data gram [%s]\n ", clientBuffer[:size])
		// 不是socks5连接请求，直接Close
		client.Close()
		return
	}
	// Connect to destiny server
	go sp.connectToServer(client)
}
func (sp *SocksProxy5Server) connectToServer(client net.Conn) {
	clientBuffer := make([]byte, sp.BufferSize)
	size, err := client.Read(clientBuffer)
	if err != nil {
		log.Printf("unable to read client Request [%s]\n ", err)
		return
	}
	// connection info from client
	// 05 01 00 01 + 目的地址(4字节） + 目的端口（2字节）
	connectionInfo := clientBuffer[:size]
	// connectionInfo[3]: IP-0x01, domain->0x03, ipv6->0x04
	if !(bytes.Equal([]byte{0x05, 0x01, 0x00, 0x01}, connectionInfo[:4]) || bytes.Equal([]byte{0x05, 0x01, 0x00, 0x03}, connectionInfo[:4])) {
		log.Printf("not proper socks5 request [%s]\n", connectionInfo)
		return
	}
	go sp.handleConnectToServer(client, connectionInfo)

}

func (sp *SocksProxy5Server) handleConnectToServer(client net.Conn, connectionInfo []byte) {

	dstHost, err := sp.parseHost(connectionInfo)
	if err != nil {
		log.Printf("unable to parse host from  gram [%s]\n", err)
		return
	}
	log.Printf("Connect [%s]\n ", dstHost)
	server, err := sp.dialer.Dial("tcp", dstHost)
	if err != nil {
		log.Printf("unable to connect to remote server [%s]\n", err)
		return
	}
	// 连接远程成功，向客户端做出回应，这里只有Connect方法
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	defer func() {
		server.Close()
		client.Close()
		log.Printf("close coonection [ %s <-> %s]\n", client.RemoteAddr(), server.RemoteAddr())
	}()
	go io.Copy(server, client)
	io.Copy(client, server)
}

func (sp *SocksProxy5Server) parseHost(connectionInfo []byte) (string, error) {
	hostType := connectionInfo[03]
	var dstPort, dstAddr string
	if hostType == 0x01 {
		dstAddr = net.IPv4(connectionInfo[4], connectionInfo[5], connectionInfo[6], connectionInfo[7]).String()
	} else if hostType == 0x03 {
		// connectionInfo[4] is the length of domain
		dstAddr = string(connectionInfo[5 : len(connectionInfo)-2])
	} else {
		return "", errors.New("ipv6 is not supported yet")
	}
	dstPort = strconv.Itoa(int(binary.BigEndian.Uint16(connectionInfo[len(connectionInfo)-2:])))
	return net.JoinHostPort(dstAddr, dstPort), nil

}
