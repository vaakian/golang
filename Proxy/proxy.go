package Proxy

import (
	"io"
	"log"
	"net"
	"net/http"
	url2 "net/url"
)

type HttpProxy struct {
	// 写一些初始化配置
	Cfg *url2.URL
}

// default listen on Cfg host
func (pxy *HttpProxy) Listen(host interface{}) {
	if host == nil || host.(string) == "" {
		http.ListenAndServe(pxy.Cfg.Host, pxy)
	} else {
		http.ListenAndServe(host.(string), pxy)
	}
}

func (pxy *HttpProxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "CONNECT" {
		log.Println("req https: " + request.Host + " start")
		ProxyHttps(writer, request)
	} else {
		log.Println("req http: " + request.Host + " start")
		ProxyHttp(writer, request)
	}
}
func ProxyHttps(writer http.ResponseWriter, request *http.Request) {
	host := request.Host
	hjk, ok := writer.(http.Hijacker)
	if !ok {
		log.Println("server doesn't support hijacking")
	}
	client, _, err := hjk.Hijack()
	if err != nil {
		return
	}

	// 连接远程
	server, err := net.Dial("tcp", host)
	if err != nil {
		return
	}
	client.Write([]byte("HTTP/1.0 200 Connection Established\r\n\r\n"))

	// 直通双向复制

	go io.Copy(server, client)

	io.Copy(client, server)

}
func ProxyHttp(writer http.ResponseWriter, request *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest(request.Method, request.URL.String(), nil)
	// 复制request header
	for k, vs := range request.Header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	// 代理请求
	res, err := client.Do(req)
	HandleErr(err)
	defer res.Body.Close()
	// 复制response header
	for k, vs := range res.Header {
		for _, v := range vs {
			writer.Header().Add(k, v)
		}
	}
	// 复制StatusCode
	writer.WriteHeader(res.StatusCode)
	// 复制body
	log.Println("req: " + request.Host + " done")
	io.Copy(writer, res.Body)
}
func HandleErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
