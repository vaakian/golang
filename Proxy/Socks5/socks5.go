package Socks5

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type SocksProxy5Server struct {
	Host    string
	Timeout time.Duration
	OffSet  int
}

func NewSocksProxy5Server(host string, timeout time.Duration) *SocksProxy5Server {
	return &SocksProxy5Server{Host: host, Timeout: timeout}
}

func (sp *SocksProxy5Server) Listen() {
	// 服务器监听并把客户端请求交由handleClientRequest处理
	srvLn, err := net.Listen("tcp", sp.Host)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("socks5 proxy server listen on: ", sp.Host)
	for {
		client, err := srvLn.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go sp.handleClientRequest(client)
	}
}
func (sp *SocksProxy5Server) handleClientRequest(client net.Conn) {
	clientBuffer := make([]byte, 1024)
	size, err := client.Read(clientBuffer)
	defer client.Close()
	if err != nil {
		log.Println("unable to read schema -> ", err)
		return
	}
	//clientBuffer[:size], []byte{0x05, bytesOfMethod, method}
	sp.OffSet = int(clientBuffer[01]) - 1
	if clientBuffer[0] == 0x05 {
		client.Write([]byte{0x05, 0x00})
	} else {
		log.Println(" unsupported data gram -> ", clientBuffer[:size])
		client.Close()
		return
	}
	// Connect to destiny server
	sp.connectToServer(client)
}
func (sp *SocksProxy5Server) connectToServer(client net.Conn) {
	clientBuffer := make([]byte, 1024)
	size, err := client.Read(clientBuffer)
	if err != nil {
		log.Println("unable to read client Request -> ", err)
		return
	}
	// connection info from client
	// 05 01 00 01 + 目的地址(4字节） + 目的端口（2字节）
	connectionInfo := clientBuffer[:size]
	// connectionInfo[3]: IP-0x01, domain->0x03, ipv6->0x04
	if !(bytes.Equal([]byte{0x05, 0x01, 0x00, 0x01}, connectionInfo[:4]) || bytes.Equal([]byte{0x05, 0x01, 0x00, 0x03}, connectionInfo[:4])) {
		log.Println("not proper socks5 data")
		log.Println(" socks5 bytes:  ", connectionInfo[:4])
		return
	}
	sp.handleConnectToServer(client, connectionInfo)

}

func (sp *SocksProxy5Server) handleConnectToServer(client net.Conn, connectionInfo []byte) {

	//dialer := net.Dialer{Timeout: sp.Timeout}
	//server, err := dialer.Dial("tcp", remoteAddr)
	defer client.Close()
	dstHost, err := sp.parseHost(connectionInfo)
	if err != nil {
		log.Println("unable to parse host from  gram -> ", err)
		return
	}
	fmt.Println("dst host: ", dstHost)
	server, err := net.Dial("tcp", dstHost)
	if err != nil {
		log.Println("unable to connect to remote server ->", err)
		return
	}
	// 连接远程成功，向客户端做出回应
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	defer server.Close()
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
