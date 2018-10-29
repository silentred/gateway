package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/silentred/gateway/route"
)

// DirectorFunc change the request to target server
type DirectorFunc func(*http.Request)

type ReverseProxyPool struct {
	pool    sync.Pool
	bufPool httputil.BufferPool
	server  *ProxyServer
}

func NewReverseProxyPool(ps *ProxyServer) *ReverseProxyPool {
	proxyPool := &ReverseProxyPool{
		server: ps,
		pool: sync.Pool{
			New: func() interface{} {
				return &httputil.ReverseProxy{
					FlushInterval: time.Second,
				}
			},
		},
		bufPool: NewBufferPool(ps.Config.Proxy.BufferSize),
	}

	return proxyPool
}

func (rpp *ReverseProxyPool) Get(s *route.Service, rw http.ResponseWriter) *httputil.ReverseProxy {
	item := rpp.pool.Get()
	if rp, ok := item.(*httputil.ReverseProxy); ok {
		rp.Director = s.Director
		rp.FlushInterval = rpp.server.Config.Proxy.FlushInterval
		rp.Transport = NewErrTransport(rpp.server.Transport)
		rp.BufferPool = rpp.bufPool
		rp.ModifyResponse = s.Modify
		rp.ErrorLog = log.New(NewErrWriter(rw), "", 0)
		return rp
	}

	return nil
}

func (rpp *ReverseProxyPool) Put(rp *httputil.ReverseProxy) {
	rpp.pool.Put(rp)
}

// deprecated
func makeDirectorFunc(target *url.URL) DirectorFunc {
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		req.URL.RawQuery = target.RawQuery
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "entree")
		}
	}
}
