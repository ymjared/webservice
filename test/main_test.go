package test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sync"
	"testing"
	"time"
)

//数组是值传递，对参数数组操作不会改变原数组的值
//切片是对数组的引用,对切片操作会修改原数组的值
func Test1(t *testing.T) {
	ta := [3]int{0, 1, 2}
	ts := ta[:]

	changeArray(ta)
	fmt.Println("ta:", ta)

	changeSlice(ts)
	fmt.Println("ta", ta)
	fmt.Println("ts:", ts)
}

func changeArray(data [3]int) {
	data[0] = 1111
	return
}

func changeSlice(data []int) {
	data[0] = 1111
	return
}

//go build -race
//写多读少使用sync.Mutex
//写少读多使用sync.RWMutex
//最大可能的读者数为实际的逻辑CPU数
func Test2(t *testing.T) {
	temp := make(map[int]int, 5)
	temp[1] = 1
	temp[2] = 1
	temp[3] = 1
	temp[10] = 1

	for i := 0; i < 10; i++ {
		go func(temp map[int]int) {
			temp[10] = 3
		}(temp)
	}

	fmt.Println(temp)
}

func Test3(t *testing.T) {
	logs := log.New(os.Stdout, "webService: ", log.Ldate|log.Ltime|log.Lshortfile)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i += 2 {
		logs.Println(rand.Intn(100))
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i += 2 {
		logs.Println(rand.Intn(100))
	}
}

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 60 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 60 * time.Second,
	MaxIdleConnsPerHost: 10,
	MaxIdleConns:        200,
}
var netCookiejar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
var netClient = &http.Client{
	Timeout:   time.Second * 60,
	Transport: netTransport,
	Jar:       netCookiejar,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func httpGet(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func Test4(t *testing.T) {
	url := "http://127.0.0.1:8889/"
	wait := sync.WaitGroup{}
	wait.Add(2000)
	for i := 0; i < 2000; i++ {
		httpGet(url)
		wait.Done()
	}
	wait.Wait()
}

//sync.cond 是协调多个等待的goroutine，
//当满足条件时可以通知一个或者多个等待的goroutine唤醒
func Test5(t *testing.T) {
	condition := sync.NewCond(&sync.Mutex{})
	wait := sync.WaitGroup{}
	wait.Add(1)

	for i := 0; i < 3; i++ {
		go func(i int) {
			log.Println("start", i)
			condition.L.Lock()
			if i == 0 {
				wait.Done()
			}
			condition.Wait()
			log.Println(i)
			time.Sleep(time.Second * 1)
			condition.L.Unlock()
		}(i)
	}

	wait.Wait()
	condition.Signal()
	time.Sleep(1 * time.Second)
	condition.Signal()
	time.Sleep(1 * time.Second)
	condition.Broadcast()
	time.Sleep(1 * time.Second)
	select {}
}

func Test6(t *testing.T) {
	//central := "http://127.0.0.1:9088/preview/central.dat"
	//fmt.Println(strings.Replace(central, "central.", "repo.", 1))

	//commKeysEncrypted := make(map[string]string)
	//if commKeysEncrypted == nil {
	//	fmt.Println("is nil")
	//} else {
	//	fmt.Println("not nil")
	//}
	bf := bytes.NewBufferString("hello" + "|")
	bf.Write([]byte("jared"))
	bf.Write([]byte("|"))
	fmt.Println(bf.String())
	fmt.Println(bf.Len())
	bf.Reset()
	fmt.Println(bf.String())
}

func Test7(t *testing.T) {
	var wg sync.WaitGroup
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()
		var buffer bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buffer, "%c", b)
		}
		fmt.Println(buffer.String())
	}

	data := []byte("golang")
	wg.Add(2)
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])
	wg.Wait()
}

func Test8(t *testing.T) {
	central, err := ioutil.ReadFile("../central.json")
	if err != nil {
		log.Println(err)
	}

	app_info, _, _, err := jsonparser.Get(central, "app_info")
	if err != nil {
		log.Println(err)
	}
	log.Printf("app_info:%s", app_info)

	app_name, _, _, err := jsonparser.Get(central, "app_name")
	if err != nil {
		log.Println(err)
	}
	log.Printf("app_name:%s", app_name)

	hosts, _, _, err := jsonparser.Get(central, "hosts")
	if err != nil {
		log.Println(err)
	}
	log.Printf("hosts:%s", hosts)
}

func Test9(t *testing.T) {
	fmt.Println("test9")
}
