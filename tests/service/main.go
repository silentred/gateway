package main

import (
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	var forever chan int
	var addrs = []string{":8085", ":8086", ":8087"}
	var prefix = []string{"hello", "world", "jason"}
	for i, item := range addrs {
		go startService(item, prefix[i])
	}
	<-forever
}

func startService(listen, prefix string) {
	e := echo.New()
	e.Any(prefix+"/*", func(ctx echo.Context) error {
		var code = 200
		randNum := rand.Int63n(100)
		// 10% returns 500 code
		if randNum >= 90 {
			code = 500
		}
		time.Sleep(time.Duration(randNum/2) * time.Millisecond)
		return ctx.String(code, ctx.Request().URL.String())
	})
	e.Start(listen)
}
