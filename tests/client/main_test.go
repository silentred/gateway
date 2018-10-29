package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"
)

func TestReq(t *testing.T) {
	//curl -i http://localhost:18088/course/coursedetail/getdetail -H 'Host: course.iget.didatrip.com' -d 'user_id=199&course_id=4'
	appKey = "ios-4.1"
	secret = "test123"
	r := getValidRequest(http.MethodPost, "entree.dev.igetget.com", "/course/course/detail", fmt.Sprint(time.Now().Unix()), "random-string")
	r.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"h":{"u": 199}, "c":{"course_id": 4} }`)))

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(dump))

	/*
		curl -XPOST 'http://localhost:18088/course/coursedetail/getdetail'
		-H 'Host: course.iget.didatrip.com' \
		-H 'Content-Type: application/x-www-form-urlencoded' \
		-H 'X-App-Key: ios-4.1' \
		-H 'X-Nonce: random-string' \
		-H 'X-Sign: palI9MztvfB2rriYQ1BmnA==' \
		-H 'X-Timestamp: 1497865388' \
		-d 'user_id=199&course_id=4'
	*/
}

func TestMd5(t *testing.T) {
	hashStr := `android
POST
application/json; charset=UTF-8
/course/course/detail

1498633544
0ec6eae7023feffb
test123`
	bytes := md5.Sum([]byte(hashStr))
	s := fmt.Sprintf("%x", bytes)
	fmt.Println(s)
	fmt.Println(base64.StdEncoding.EncodeToString([]byte(s)))
}
