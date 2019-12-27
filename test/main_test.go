package test

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"
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

//var netCookiejar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
var netClient = &http.Client{
	Timeout:   time.Second * 60,
	Transport: netTransport,
	//Jar:       netCookiejar,
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
	registry, err := ioutil.ReadFile("../registry.json")
	if err != nil {
		log.Println(err)
	}
	jsonparser.ArrayEach(registry, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		certdata, err := jsonparser.GetString(value, "cert")
		if err != nil {
			log.Println("Mutual Auth/2-Way SSL: Failed to parse client certificate: " + err.Error())
			return
		}
		keydata, err := jsonparser.GetString(value, "key")
		if err != nil {
			log.Println("Mutual Auth/2-Way SSL: Failed to parse client private key: " + err.Error())
			return
		}
		fmt.Println("cert:", certdata)
		fmt.Println("key:", keydata)
	}, "tls_client", "certs")
}

func Test10(t *testing.T) {
	central := "http://127.0.0.1:9088/preview/central.dat"
	repoconf := "./main/repo.dat"
	repourlconf := "./"

	centralurl, _ := url.Parse(central)
	repourl, _ := url.Parse(repoconf)
	repodatUrl := centralurl.ResolveReference(repourl)
	fmt.Println("repo.dat: ", repodatUrl.String())

	repourlconfurl, _ := url.Parse(repourlconf)
	fmt.Println("repo resource: ", repodatUrl.ResolveReference(repourlconfurl).String())

}

func Test11(t *testing.T) {
	rule := `\.do\\?`
	examples := []string{`cockpit/vendor/deryrr/framework7/cssframework7.ios.min.css`}

	r, err := regexp.Compile(rule)
	if err != nil {
		log.Fatal(err)
	}
	for _, example := range examples {
		if r.MatchString(example) {
			fmt.Println(example)
		}
	}
}

func Test12(t *testing.T) {
	path := "weex/financeProducts.weex.js"
	salt := "csii-fn-key"
	target := fmt.Sprintf("%x.dat", sha256.Sum256([]byte(salt+path+salt)))
	if target == "f6b518e915b775490846fe952bf78272fbd83c95933d2356446221467991b068.dat" {
		fmt.Println("ok")
	}
}

func Test13(t *testing.T) {
	dbpath := "./ss.dat"
	key := "_ENV"
	dat, err := ioutil.ReadFile(dbpath)
	if err != nil {
		log.Println("Secure Store: Getting key " + key + err.Error())
		return
	}
	str, err := jsonparser.GetString(dat, key)
	if err != nil {
		log.Println("Secure Store: Getting key " + key + err.Error())
		return
	}
	fmt.Println(str)
}

func Test14(t *testing.T) {
	http.Handle("/", http.FileServer(http.Dir("../")))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type Weapon int

const (
	shouqiang Weapon = iota
	dapao
	feiji
	huojian
)

func Test15(t *testing.T) {
	fmt.Println(shouqiang)
	fmt.Println(dapao)
	fmt.Println(feiji)
	fmt.Println(huojian)
}

func (w Weapon) String() string {
	switch w {
	case shouqiang:
		return "shouqiang"
	case dapao:
		return "dapao"
	case feiji:
		return "feiji"
	case huojian:
		return "huojian"
	}
	return ""
}

func Test16(t *testing.T) {
	src := make([]int, 10)
	copySrc := make([]int, 10)

	for i := 0; i < 10; i++ {
		src[i] = i + 1
	}
	//引用类型
	refSrc := src
	//复制
	copy(copySrc, src)
	src[0] = 100
	fmt.Println("src:", src)
	fmt.Println("refSrc:", refSrc)
	fmt.Println("copySrc:", copySrc)

	match, err := filepath.Glob("./*")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(match)

	baseName := filepath.Base("/a/b/c/e.txt")
	fmt.Println(baseName)

	e := filepath.Ext("/Users/1.txt")
	fmt.Println(e)
	s := "i i"
	fmt.Println(utf8.RuneCountInString(s))

	notALetter := func(char rune) bool {
		return unicode.IsLetter(char)
	}
	fmt.Println(len(strings.FieldsFunc(s, notALetter)), strings.FieldsFunc(s, notALetter))

	gap := 10 + 5 - len("Word") - len("Frequency")
	fmt.Printf("Word %*s%s\n", gap, " ", "Frequency")
}

type Invoke interface {
	Call(interface{})
}

//定义函数类型
type FuncCaller func(interface{})

func (f FuncCaller) Call(p interface{}) {
	f(p)
}

//回调函数
func Test17(t *testing.T) {
	var data int
	var funcCaller Invoke
	callBack(10, func(i interface{}) {
		data = i.(int)
	})
	fmt.Println("callback: ", data)

	funcCaller = FuncCaller(func(i interface{}) {
		data = i.(int)
	})

	funcCaller.Call(20)
	fmt.Println("funcCall: ", data)
}

func callBack(d int, f FuncCaller) {
	f(d)
}

//内存缓存
const LIM = 41

var fibs [LIM]uint64

func Test19(t *testing.T) {
	var result uint64 = 0
	start := time.Now()
	for i := 1; i < LIM; i++ {
		result = fibonacci3(i)
		fmt.Printf("数列第 %d 位: %d\n", i, result)
	}
	end := time.Now()
	delta := end.Sub(start)
	fmt.Printf("程序的执行时间为: %s\n", delta)
}
func fibonacci3(n int) (res uint64) {
	// 记忆化：检查数组中是否已知斐波那契（n）
	if fibs[n] != 0 {
		res = fibs[n]
		return
	}
	if n <= 2 {
		res = 1
	} else {
		res = fibonacci3(n-1) + fibonacci3(n-2)
	}
	fibs[n] = res
	return
}

func Test20(t *testing.T) {
	data := []byte("golang")
	br := bytes.NewReader(data)
	bufr := bufio.NewReader(br)
	var target [256]byte
	n, err := bufr.Read(target[:])
	fmt.Println(string(target[:n]), n, err)
}

type Driver interface {
	Run()
}
type Driver1 interface {
	Run()
}

type Car struct {
	Name string
}

func (c Car) Run() {
	fmt.Println("Car: Run")
}

func Test21(t *testing.T) {
	var ei interface{}
	ei = Car{}
	di := ei.(Driver)
	di.Run()
	di = ei.(Driver1)
	di.Run()
}
