package Proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	_ "io"
	"log"
	"net"
	"net/http"
	url2 "net/url"
	"testing"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// echo server
func Test_socket(t *testing.T) {
	ln, err := net.Listen("tcp", ":1080")
	HandleErr(err)
	for {
		conn, err := ln.Accept()
		HandleErr(err)
		go func(conn net.Conn) {
			defer conn.Close()
			buf := make([]byte, 500)
			size, err := conn.Read(buf)
			HandleErr(err)
			strBuf := string(buf[:size])
			fmt.Println(strBuf)
			fmt.Println("------------------------------------")
			conn.Write(buf[:size])
		}(conn)
	}
}

func Test_websocket(t *testing.T) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := Upgrader.Upgrade(w, r, nil)
		HandleErr(err)
		defer conn.Close()

		for {
			user := &User{}
			err := conn.ReadJSON(user)
			HandleErr(err)
			log.Print(user)
			for i := 0; i < 10; i++ {
				user.Age++
				err = conn.WriteJSON(user)
				HandleErr(err)
			}
		}
	})
	http.ListenAndServe(":8080", nil)
}

// http proxy
func Test_httpProxy(t *testing.T) {
	url := "http://:1080"
	Cfg, err := url2.Parse(url)
	HandleErr(err)
	pxy := &HttpProxy{Cfg: Cfg}
	// nil 默认
	pxy.Listen(nil)
}

func Test_httpsRequest(t *testing.T) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", "https://baidu.com", nil)
	HandleErr(err)
	res, err := cli.Do(req)
	HandleErr(err)

	log.Println(res)
}
