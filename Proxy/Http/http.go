package Http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

type HttpProxyServer struct {
	Host string
}

func NewHttpProxyServer(host string) *HttpProxyServer {
	return &HttpProxyServer{Host: host}
}

func (htp *HttpProxyServer) Listen() {
	ln, err := net.Listen("tcp", htp.Host)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("http proxy server listen on: " + htp.Host)
	for {
		client, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go htp.handleClientRequest(client)
	}
}
func (htp *HttpProxyServer) handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()
	// 用来存放客户端数据的缓冲区
	var b [2048]byte
	//从客户端获取数据
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}

	var method, URL, address string
	// 从客户端数据读入method，url
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &URL)
	hostPortURL, err := url.Parse(URL)
	if err != nil {
		log.Println(err)
		return
	}
	address = URL
	if method != "CONNECT" {
		//否则为http协议
		address = hostPortURL.Host
		// 如果host不带端口，则默认为80
		if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
			address = hostPortURL.Host + ":80"
		}
	}
	log.Printf("%s -> %s -> %s\n", client.RemoteAddr(), client.LocalAddr(), address)
	//获得了请求的host和port，向服务端发起tcp连接
	htp.connectToServer(client, address, method, b[:n])
}
func (htp *HttpProxyServer) connectToServer(client net.Conn, address string, method string, ctx []byte) {
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	//如果使用https协议，需先向客户端表示连接建立完毕，https直接tcp转发即可
	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else { //如果使用http协议，需将从客户端得到的http请求信息转发给服务端
		server.Write(ctx)
	}
	//将客户端的请求转发至服务端，将服务端的响应转发给客户端。io.Copy为阻塞函数，文件描述符不关闭就不停止
	go io.Copy(server, client)
	io.Copy(client, server)
	// http: 写请求体到tcp连接，tcp返回，并关闭连接。
	// https: 不写请求体，服务器443端口建立tcp连接，双方通过该连接传递相应信息。
}
