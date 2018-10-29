package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/silentred/glog"
	"github.com/silentred/gateway/admin"
	"github.com/silentred/gateway/config"
	"github.com/silentred/gateway/config/consul"
	"github.com/silentred/gateway/config/etcd"
	"github.com/silentred/gateway/guard"
	"github.com/silentred/gateway/guard/rate_limiter"
	"github.com/silentred/gateway/guard/replay"
	"github.com/silentred/gateway/metric"
	"github.com/silentred/gateway/proxy"
	"github.com/silentred/gateway/reactor"
	"github.com/silentred/gateway/reactor/circuit_breaker"
)

const (
	version = "1.0.2"
)

var (
	cfgFile string
	showVer bool
	GitHash = "None"
	BuildTS = "None"
)

func init() {
	flag.StringVar(&cfgFile, "c", "config.yaml", "path of config file")
	flag.BoolVar(&showVer, "v", false, "show version")
}

func usage() {
	fmt.Fprintf(os.Stderr, `
 _____       _
| ____|_ __ | |_ _ __ ___  ___
|  _| | '_ \| __| '__/ _ \/ _ \
| |___| | | | |_| | |  __/  __/
|_____|_| |_|\__|_|  \___|\___|

Version: %s 
BuildTS: %s
GitHash: %s

`, version, BuildTS, GitHash)

	flag.PrintDefaults()
}

func main() {
	var err error

	flag.Parse()
	flag.Usage = usage
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	if showVer {
		usage()
		return
	}

	// load config
	cfg, err := config.Load(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	// register guards
	initGuards(cfg)
	// register reactors
	initReactors(cfg)

	ps := proxy.NewProxyServer(cfg)

	// TODO: file backend
	var backend config.Backend
	switch cfg.Route.Backend {
	case "etcd":
		backend = etcd.NewBackend(ps.Table, &cfg.Etcd)
	case "consul":
		backend = consul.NewBackend(ps.Table, &cfg.Consul)
	default:
		log.Fatalf("not support backend: [%s]", cfg.Route.Backend)
	}

	go backend.Watch()

	// start metrics
	go metric.Start("/metrics", ":7088")
	go admin.Start(cfg, ps.Table, backend)
	go admin.StartWebUI(cfg)
	go listenAndServe(cfg, ps)

	handleSignal()
}

func listenAndServe(cfg *config.Config, h http.Handler) {
	var err error
	// same as http,ListenAndServe(addr, handler); to control the timeouts
	server := http.Server{
		Addr:         cfg.Listen.Address,
		Handler:      h,
		ReadTimeout:  cfg.Listen.ReadTimeout,
		WriteTimeout: cfg.Listen.WriteTimeout,
		IdleTimeout:  cfg.Listen.IdleTimeout,
	}

	// use for-loop and sleep to ensure this app does not quit.
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("listen error: %v \n", err)
	}
}

func handleSignal() {
	var (
		c chan os.Signal
		s os.Signal
	)
	c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,
		syscall.SIGINT, syscall.SIGSTOP, syscall.SIGPIPE, syscall.SIGCHLD)
	// Block until a signal is received.
	for {
		s = <-c
		glog.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			// TODO close
			return
		case syscall.SIGHUP, syscall.SIGPIPE, syscall.SIGCHLD:
			// TODO reload
			//return
		default:
			return
		}
	}
}

func initGuards(cfg *config.Config) {
	var err error

	sg := guard.NewSignGuard(cfg.Guard.Secret)
	now := func() time.Time {
		return time.Now()
	}
	tg := guard.NewTimeGuard(now, int64(cfg.Guard.TimeRange))

	guard.Register("sign", sg)
	guard.Register("time", tg)
	registerRateLimiter(cfg)
	registerReplay(cfg)

	guard.DefaultGroup, err = guard.GetGroup("default", "sign", "time", "rate_limiter", "replay")
	if err != nil {
		//Log
		log.Fatal(err)
	}
}

func registerRateLimiter(cfg *config.Config) {
	var ttl time.Duration
	var err error
	ttl, err = time.ParseDuration(cfg.Reactor.RateLimiter.TTL)
	if err != nil {
		log.Fatal(err)
	}

	ttlKV := circuit_breaker.NewMemTTLKV(ttl)
	attlKV := rate_limiter.NewMemAtomicTTLKV(*ttlKV, cfg.Reactor.RateLimiter.LockNum)
	rl := rate_limiter.NewRateLimiter(attlKV, cfg.Reactor.RateLimiter.ThresholdCount)
	guard.Register("rate_limiter", rl)
}

func registerReplay(cfg *config.Config) {
	var ttl time.Duration
	var err error
	ttl, err = time.ParseDuration(cfg.Reactor.RateLimiter.TTL)
	if err != nil {
		log.Fatal(err)
	}

	c := circuit_breaker.NewMemTTLKV(ttl)
	rr := replay.NewReplayReactor(c, ttl)
	guard.Register("replay", rr)
}

func initReactors(cfg *config.Config) {
	var err error

	registerCircuitBreaker(cfg)

	reactor.DefaultGroup, err = reactor.GetGroup("defualt", "cb")
	if err != nil {
		//TODO log
		log.Fatal(err)
	}
}

func registerCircuitBreaker(cfg *config.Config) {
	var thresholdDuration time.Duration
	var blockDuration time.Duration
	var err error

	thresholdDuration, err = time.ParseDuration(cfg.Reactor.CircuitBreaker.ThresholdDurationStr)
	if err != nil {
		log.Fatal(err)
	}
	blockDuration, err = time.ParseDuration(cfg.Reactor.CircuitBreaker.BlockDurationStr)
	if err != nil {
		log.Fatal(err)
	}

	if thresholdDuration < time.Minute {
		glog.Infof("[init CB] threshold duration too small, current=%s. minimal value is 1 minute", thresholdDuration)
		thresholdDuration = time.Minute
	}
	if blockDuration < time.Minute {
		glog.Infof("[init CB] block duration too small, current=%s. minimal value is 1 minute", blockDuration)
		blockDuration = time.Minute
	}

	c := circuit_breaker.NewMemTTLKV(thresholdDuration)
	s := circuit_breaker.NewBinaryStrategy()
	cb := circuit_breaker.NewCircuitBreaker(c, s, cfg.Reactor.CircuitBreaker.ThresholdCount, thresholdDuration, blockDuration)
	reactor.Register("cb", cb)
}
