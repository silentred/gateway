package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/silentred/gateway/config"
	"github.com/silentred/gateway/route"
)

func TestProxyServer(t *testing.T) {
	var err error
	cfg := config.DefaultConfig()
	ps := NewProxyServer(&cfg)

	// http
	rt := route.NewRoute("www.fakeweb.com", "/v1/hello")
	target := route.NewTarget("", "localhost:7000", 1)
	service := route.NewService("http-service", "", "/v1", target, nil, nil)
	assert.NotNil(t, service.Targets)
	ps.Table.Add(rt, service)

	// test Connected number
	for i := 0; i < 10; i++ {
		ps.incrConn()
	}
	assert.Equal(t, uint64(10), ps.Connected)
	ps.decrConn()
	assert.Equal(t, uint64(9), ps.Connected)

	// start backend server
	go startHTTPServer()

	req, err := http.NewRequest(http.MethodGet, "http://www.fakeweb.com/v1/hello?a=b&c=d", nil)
	assert.NoError(t, err)
	reqNotFound, _ := http.NewRequest(http.MethodGet, "http://noexist.com/v1/hello?a=b&c=d", nil)

	tbl := []struct {
		req        *http.Request
		bodyPrefix string
		code       int
	}{
		{req, "GET /hello?a=b&c=d", 200},
		{reqNotFound, `{"code":100404`, 404},
	}

	for idx, item := range tbl {
		var record = httptest.NewRecorder()
		// handle request
		ps.ServeHTTP(record, item.req)

		assert.Equal(t, item.code, record.Code)
		body := string(record.Body.Bytes())
		assert.True(t, strings.HasPrefix(body, item.bodyPrefix))
		t.Logf("idx:%d body: %s", idx, body)
	}

}

func startHTTPServer() {
	// http
	http.HandleFunc("/hello", func(rw http.ResponseWriter, r *http.Request) {
		//time.Sleep(5 * time.Second)
		b, _ := httputil.DumpRequest(r, true)
		fmt.Fprint(rw, string(b))
	})

	// websocket
	var upgrader = websocket.Upgrader{}
	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			log.Println(err)
		}
		for {
			// read message
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
					fmt.Println(err)
				}
				break
			}
			log.Println(msg)
		}
	})

	http.ListenAndServe(":7000", nil)
}
