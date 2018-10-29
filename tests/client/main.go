package main

import (
	"crypto/md5"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"strings"

	"sync"

	"github.com/fatih/color"
	"github.com/silentred/gateway/util"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	addr   string
	secret string
	appKey string
	loop   bool
	client = util.NewHTTPClient(10, nil)

	// color sprintf
	// error
	red sprintfFunc
	// highlight
	blue sprintfFunc
	// good
	green sprintfFunc
	// warning
	yellow sprintfFunc
)

type sprintfFunc func(format string, a ...interface{}) string

func main() {
	setup()
	// test concurrent visits; test sign;
	for loop {
		testSign()
		time.Sleep(5 * time.Second)
	}
	// test concurrent visits; test sign;
	testSign()
	// test timestamp
	testTimestamp()
	// test CB
	testCircuitBreaker()
	// test rate limiter
	testRateLimiter()
	// test replay
	testReplay()
}

func setup() {
	// parse flag
	flag.StringVar(&addr, "addr", "localhost:8088", "proxy address")
	flag.StringVar(&secret, "secret", "test123", "proxy secret")
	flag.StringVar(&appKey, "appKey", "ios-client", "client key")
	flag.BoolVar(&loop, "loop", false, "keep sending requests to entree")
	flag.Parse()

	// log
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	// color
	red = color.New(color.FgRed).SprintfFunc()
	blue = color.New(color.FgCyan).SprintfFunc()
	green = color.New(color.FgGreen).SprintfFunc()
	yellow = color.New(color.FgYellow).SprintfFunc()

	// log.Println("Showing the meaning of each color:")
	// log.Println(red("Error"))
	// log.Println(yellow("Warning"))
	// log.Println(blue("Highlight"))
	// log.Println(green("Good"))

	// start backend service
}

// RandomString returns the random string with length of n
func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func sign(key, method, conType, path, query, time, nonce, secret string) string {
	hashStr := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s", key, method, conType, path, query, time, nonce, secret)
	bytes := md5.Sum(util.Slice(hashStr))
	md5Str := fmt.Sprintf("%x", util.String(bytes[:]))
	return base64.StdEncoding.EncodeToString(util.Slice(md5Str))
}

func getValidRequest(method, host, path, time, nonce string) *http.Request {
	fullURL := fmt.Sprintf("http://%s%s", addr, path)
	u, err := url.Parse(fullURL)
	if err != nil {
		log.Println(red("parseURL: %v", err))
		return nil
	}
	queryStr := u.Query().Encode()
	s := sign(appKey, method, "application/x-www-form-urlencoded", u.Path, queryStr, time, nonce, secret)

	header := map[string]string{
		"X-App-Key":    appKey,
		"Content-Type": "application/x-www-form-urlencoded",
		"X-Timestamp":  time,
		"X-Nonce":      nonce,
		"X-Sign":       s,
	}

	req, _ := util.NewHTTPReqeust(method, fullURL, nil, header, nil)
	// Note:  prevents the connection from being re-used.
	//  request.Close [a bool] indicates whether to close the connection after replying to this request (for servers) or after sending the request (for clients)
	req.Close = true
	req.Host = host
	return req
}

func validateReq(r *http.Request, bodySubStr string) error {
	body, code, err := client.Do(r)
	if err != nil || !strings.Contains(string(body), bodySubStr) {
		return fmt.Errorf(red("bad routing: path=%s body=%s code=%d err=%v", r.URL.String(), string(body), code, err))
	}
	return nil
}

func validateBadReq(r *http.Request, code int) error {
	body, c, err := client.Do(r)
	if c < code {
		return fmt.Errorf(red("bad req: code:%d path:%s body:%s err:%v", c, r.URL.String(), string(body), err))
	}
	return nil
}

// test #1
func testSign() {
	var wg sync.WaitGroup
	var lastError, err error

	log.Println(blue("--- Case: 签名验证 ---"))
	log.Println(blue("hello, world service 共500次 并发"))

	// b, _ := httputil.DumpRequest(r, false)
	// log.Println(string(b))
	for idx := 0; idx < 5; idx++ {
		for i := 0; i < 50; i++ {
			r := getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/hello-service?z=y&a=b&c=d", fmt.Sprint(time.Now().Unix()), randomString(16))
			r2 := getValidRequest(http.MethodGet, "world.luoji.com", "/world/world-service?z=y&a=b&c=d", fmt.Sprint(time.Now().Unix()), randomString(16))

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = validateReq(r, "hello-service")
				if err != nil {
					log.Println(err)
					lastError = err
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = validateReq(r2, "world-service")
				if err != nil {
					log.Println(err)
					lastError = err
				}
			}()
		}
		wg.Wait()
		log.Println(blue("%d00次完成", idx+1))
	}

	if lastError != nil {
		log.Println(lastError)
	} else {
		log.Println(green("--- 成功 ---"))
	}
}

// test #2
func testTimestamp() {
	var err error
	log.Println(blue("--- Case: Timestamp验证 ---"))
	log.Println(blue("hello service 3次 "))

	now := fmt.Sprint(time.Now().Unix())
	past := fmt.Sprint(time.Now().Add(-1 * 10 * time.Minute).Unix())
	future := fmt.Sprint(time.Now().Add(10 * time.Minute).Unix())

	var times = []string{now, past, future}

	for _, item := range times {
		r := getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/hello-service?z=y&a=b&c=d", item, randomString(16))
		err = validateReq(r, "hello-service")
		if err != nil {
			log.Println(err)
			break
		}
	}

	if err != nil {
		log.Println(err)
	} else {
		log.Println(green("--- 成功 ---"))
	}
}

// test #3
func testCircuitBreaker() {
	var err error
	var threshold = 60
	//var thresholdDuration = 30 * time.Second
	//var blockDuration = 30 * time.Second
	var timeStr = fmt.Sprint(time.Now().Unix())
	badReq := getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/bad?z=y&a=b&c=d", timeStr, randomString(16))
	goodReq := getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/hello-service?z=y&a=b&c=d", timeStr, randomString(16))

	log.Println(blue("--- Case: CircuitBreaker 验证 ---"))

	// #1
	log.Println(blue("60 bad requests in 30s, and NOT BLOCK requests "))
	var wg sync.WaitGroup

	for i := 0; i < threshold; i++ {
		badReq := getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/bad?z=y&a=b&c=d", timeStr, randomString(16))
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := validateBadReq(badReq, 500)
			if err != nil {
				log.Println(err)
			}
		}()
	}
	wg.Wait()

	err = validateReq(goodReq, "hello-service")
	if err != nil {
		log.Println(err)
	}
	log.Println(blue("OK, not block"))

	// #2
	log.Println(blue("one more bad requests in last 30s, and BLOCK 1/2 all requests"))
	var passCnt int
	var badCnt int
	badReq = getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/bad?z=y&a=b&c=d", timeStr, randomString(16))
	err = validateBadReq(badReq, 500)
	if err != nil {
		log.Println(err)
	}
	for passCnt < 61 {
		badReq := getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/bad?z=y&a=b&c=d", timeStr, randomString(16))
		_, code, _ := client.Do(badReq)
		if code == 500 {
			passCnt++
		} else {
			badCnt++
		}
	}
	log.Println(blue("pass:%d blocked:%d", passCnt, badCnt))

	// #3
	log.Println(blue("60 more bad requests in last 30s, and BLOCK 3/4 all requests"))
	passCnt = 0
	badCnt = 0
	for passCnt < 61 {
		badReq := getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/bad?z=y&a=b&c=d", timeStr, randomString(16))
		_, code, _ := client.Do(badReq)
		if code == 500 {
			passCnt++
		} else {
			badCnt++
		}
	}
	log.Println(blue("pass:%d blocked:%d", passCnt, badCnt))

	// log.Println(blue("sleep for 30s"))
	// time.Sleep(thresholdDuration)

	log.Println(blue("try a good request now"))
	goodReq = getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/hello-service?z=y&a=b&c=d", timeStr, randomString(16))
	err = validateReq(goodReq, "hello-service")
	if err != nil {
		log.Println(err)
	}

	log.Println(green("--- 成功 ---"))
}

// test #4
func testRateLimiter() {
	var err error

	log.Println(blue("--- Case: RateLimiter验证 ---"))
	log.Println(blue("hello service 60次"))

	// use X-UID as limit id
	now := fmt.Sprint(time.Now().Unix())

	// not block
	for i := 0; i < 61; i++ {
		r := getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/hello-service?z=y&a=b&c=d", now, randomString(16))
		r.Header.Add("X-UID", "100100")
		err = validateReq(r, "hello-service")
		if err != nil {
			log.Println(err)
		}
	}

	// block
	log.Println(blue("another 10 visits, block"))
	for i := 0; i < 10; i++ {
		r := getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/hello-service?z=y&a=b&c=d", now, randomString(16))
		r.Header.Add("X-UID", "100100")
		err = validateBadReq(r, 403)
		if err != nil {
			log.Println(err)
		}
	}

	// sleep 30s
	log.Println(blue("sleeping for 30s"))
	time.Sleep(30 * time.Second)
	for i := 0; i < 10; i++ {
		r := getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/hello-service?z=y&a=b&c=d", now, randomString(16))
		r.Header.Add("X-UID", "100100")
		err = validateReq(r, "hello-service")
		if err != nil {
			log.Println(err)
		}
	}

	log.Println(green("--- 成功 ---"))
}

// test #5
func testReplay() {
	var err error

	log.Println(blue("--- Case: Replay 验证 ---"))
	log.Println(blue("hello service 1次"))

	now := fmt.Sprint(time.Now().Unix())
	r := getValidRequest(http.MethodGet, "hello.luoji.com", "/hello/hello-service?z=y&a=b&c=d", now, randomString(16))
	err = validateReq(r, "hello-service")
	if err != nil {
		log.Println(err)
	}

	log.Println(blue("hello service 1次, BLOCK"))
	err = validateBadReq(r, 403)
	if err != nil {
		log.Println(err)
	}

	log.Println(blue("sleeping for 30s, not block again"))
	time.Sleep(30 * time.Second)
	err = validateReq(r, "hello-service")
	if err != nil {
		log.Println(err)
	}
	log.Println(green("--- 成功 ---"))
}
