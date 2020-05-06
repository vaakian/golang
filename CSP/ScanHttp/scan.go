package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net"
	"net/http"
	"sync"
	"time"
)

type httpRequest struct {
	url string
	res string
}

func main() {
	cdir := flag.String("cdir", "", "cdir example: 172.16.0.0/24")
	port := flag.String("port", "80", "http port")
	threadCount := flag.Int("thread", 32, "http requests thread")
	bufferSize := flag.Int("buffer", 64, "http requests job buffer size")
	flag.Parse()
	if _, _, err := net.ParseCIDR(*cdir); err != nil {
		fmt.Println("invalid cdir")
		flag.PrintDefaults()
		return
	}

	hostJob := make(chan string, *bufferSize)
	result := make(chan httpRequest, *bufferSize)
	// 创建goroutine
	wg := new(sync.WaitGroup)
	wg.Add(*threadCount)
	for i := 0; i < *threadCount; i++ {
		go worker(hostJob, result, wg)
	}
	go func() {
		wg.Wait()
		//	一旦所有work完成，则通知result完成
		close(result)
	}()
	go addJob(hostJob, *cdir, *port)
	printResult(result)
	fmt.Printf("cdir: %s\nport: %s\nthread: %d\nbuffer: %d\n", *cdir, *port, *threadCount, *bufferSize)
}
func printResult(result chan httpRequest) {
	for res := range result {
		fmt.Printf("%s --> %s\n", res.url, res.res)
	}
}
func addJob(job chan<- string, cdir string, port string) {
	ips, _ := hosts(cdir)
	for _, ip := range ips {
		host := "http://" + net.JoinHostPort(ip, port)
		job <- host
	}
	close(job)
}
func worker(hostJob <-chan string, result chan<- httpRequest, wg *sync.WaitGroup) {
	for url := range hostJob {
		res, err := httpDetect(url)
		if err == nil {
			result <- res
		}
	}
	wg.Done()
}
func httpDetect(host string) (httpRequest, error) {
	cli := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		return httpRequest{}, err
	}
	res, err := cli.Do(req)
	if err != nil {
		return httpRequest{}, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return httpRequest{}, err
	}
	title := doc.Find("title").Text()
	if title == "" {
		title = doc.Text()
		if len(title) >= 50 {
			title = title[:50]
		}
	}
	return httpRequest{url: host, res: title}, nil
}

func hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
