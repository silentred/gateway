package proxy

import (
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/silentred/glog"
	"github.com/silentred/gateway/config"
	"github.com/silentred/gateway/metric"
	"github.com/silentred/gateway/route"
)

type ProxyServer struct {
	Config    *config.Config
	Transport http.RoundTripper
	Table     *route.Table
	Connected uint64
	pool      *ReverseProxyPool
}

func NewProxyServer(cfg *config.Config) *ProxyServer {
	ps := &ProxyServer{
		Config:    cfg,
		Transport: NewTransport(cfg),
		Table:     route.NewTable(),
	}
	ps.pool = NewReverseProxyPool(ps)
	return ps
}

func NewTransport(cfg *config.Config) *http.Transport {
	return &http.Transport{
		ResponseHeaderTimeout: cfg.Proxy.ResponseHeaderTimeout,
		MaxIdleConnsPerHost:   cfg.Proxy.MaxConn,
		Dial: (&net.Dialer{
			Timeout:   cfg.Proxy.DialTimeout,
			KeepAlive: cfg.Proxy.KeepAliveTimeout,
		}).Dial,
	}
}

func (ps *ProxyServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var err, errReactor error
	var upgrade string

	// recover panic
	defer func() {
		if catch := recover(); catch != nil {
			var err error
			switch catch := catch.(type) {
			case error:
				err = catch
			default:
				err = fmt.Errorf("%v", catch)
			}
			stack := make([]byte, ps.Config.Recover.StackSize)
			length := runtime.Stack(stack, !ps.Config.Recover.DisableStackAll)
			if !ps.Config.Recover.DisablePrintStack {
				glog.Errorf("[PANIC] %s %s", err, stack[:length])
			}
		}
	}()

	s := ps.Table.FindByRequest(r)
	// 404 not found
	if s == nil {
		glog.Errorf("[Route] not found host:%s path:%s", r.Host, r.URL.Path)
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rw, route.ErrNotFound)
		return
	}

	// Guards try to Reject
	if s.Guards != nil {
		err = s.Guards.Reject(r)
	}
	if err != nil {
		glog.Errorf("[Guard] host:%s path:%s err:%s", r.Host, r.URL.Path, err)
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprint(rw, err)
		return
	}

	// copy URL
	// var orignalURL = *r.URL
	// r.URL.Path = s.StripPrefix(r.URL.Path)
	if s.Reactors != nil {
		errReactor = s.Reactors.Reject(r)
	}
	if errReactor != nil {
		glog.Errorf("[Reactor] host:%s path:%s err:%s", r.Host, r.URL.Path, errReactor)
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprint(rw, errReactor)
		return
	}
	// set back
	// r.URL = &orignalURL

	upgrade = r.Header.Get("Upgrade")
	glog.Debugf("[Upgrade] upgrade=%s svcName=%s", upgrade, s.Name)

	if strings.ToLower(upgrade) == "websocket" {
		ps.handleWebsocket(rw, r, s)
	} else {
		ps.handleHTTP(rw, r, s)
	}
}

func (ps *ProxyServer) handleWebsocket(rw http.ResponseWriter, r *http.Request, s *route.Service) {
	// TODO
	wsProxy := WebsocketProxy{Director: s.Director}
	wsProxy.ServeHTTP(rw, r)
}

func (ps *ProxyServer) handleHTTP(rw http.ResponseWriter, r *http.Request, s *route.Service) {
	var err error
	var start time.Time
	var duration time.Duration
	var resp *http.Response
	// copy original request
	var originalReq = *r
	var originalReqURL = *r.URL

	ps.incrConn()
	metric.AddConnectionNumber(1, metric.HTTPConnection)
	defer func() {
		ps.decrConn()
		metric.AddConnectionNumber(-1, metric.HTTPConnection)
	}()

	// get httpProxy
	httpProxy := ps.pool.Get(s, rw)
	defer ps.pool.Put(httpProxy)

	// record time
	start = time.Now()
	// s.Director() change the request's path
	httpProxy.ServeHTTP(rw, r)
	duration = time.Since(start)

	// Reactor may need observe error
	if et, ok := httpProxy.Transport.(*ErrTransport); ok {
		err = et.LastErr
		resp = et.Response

		// collect metrics
		metric.HTTPDurationObserve(resp, s, duration)

		// Reactor needs observe err and request
		originalReq.URL = &originalReqURL
		s.Reactors.ObserveError(&originalReq, resp, err)
		if err != nil {
			glog.Errorf("[RoundTripper] host:%s path:%s err:%s", r.Host, r.URL.Path, err)
		}
	}

	// log to access.log
	if resp != nil {
		glog.Infof("[Access] host=%s method=%s uri=%s path=%s status=%d latency_str=%s latency=%d",
			r.Host, r.Method, r.URL.RequestURI(), r.URL.Path, resp.StatusCode,
			duration.String(), duration.Nanoseconds()/1000)
	}
}

func (ps *ProxyServer) incrConn() {
	atomic.AddUint64(&ps.Connected, 1)
}

func (ps *ProxyServer) decrConn() {
	var i int64 = -1
	atomic.AddUint64(&ps.Connected, uint64(i))
}
