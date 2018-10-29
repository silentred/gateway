package proxy

import (
	"io"
	"net"
	"net/http"

	"github.com/silentred/glog"

	"fmt"

	"github.com/silentred/gateway/util"
)

var (
	ErrHijacker    = util.NewError(100550, "not hijacker")
	ErrHijack      = util.NewError(100551, "hijack failed")
	ErrDialRemote  = util.NewError(100555, "dial remote failed")
	ErrWriteRemote = util.NewError(100556, "write remote failed")
)

type WebsocketProxy struct {
	Director DirectorFunc
}

// ServeHTTP implements HTTP server
func (wsp *WebsocketProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var inConn net.Conn
	var origin string

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, ErrHijacker.Error(), http.StatusBadGateway)
		glog.Error(ErrHijacker)
		return
	}

	inConn, _, err = hj.Hijack()
	if err != nil {
		http.Error(w, ErrHijack.Error(), http.StatusBadGateway)
		glog.Error(err)
		return
	}
	defer inConn.Close()

	wsp.Director(r)
	glog.Debugf("Origin=%s TargetHost=%s", r.Header.Get("Origin"), r.Host)
	origin = fmt.Sprintf("http://%s", r.Host)
	r.Header.Set("Origin", origin)

	outConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, ErrDialRemote.Error(), http.StatusBadGateway)
		glog.Error(err)
		return
	}
	defer outConn.Close()

	err = r.Write(outConn)
	if err != nil {
		http.Error(w, ErrWriteRemote.Error(), http.StatusBadGateway)
		glog.Error(err)
		return
	}

	errc := make(chan error, 2)
	cp := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		errc <- err
	}

	go cp(outConn, inConn)
	go cp(inConn, outConn)
	err = <-errc
	if err != nil && err != io.EOF {
		//log.Printf("[INFO] WS error for %s. %s", r.URL, err)
		glog.Error(err)
	}
}
