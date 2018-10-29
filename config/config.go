package config

import (
	"os"
	"time"

	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Listen  `yaml:"listen"`
	Recover `yaml:"recover"`
	Proxy   `yaml:"proxy"`
	Route   `yaml:"route"`
	Guard   `yaml:"guard"`
	Reactor `yaml:"reactor"`
	Admin   `yaml:"admin"`
	Metric  `yaml:"metric"`
	Consul  `yaml:"consul"`
	Etcd    `yaml:"etcd"`
	WebUI   `yaml:"webui"`
}

type Listen struct {
	Address         string `yaml:"address"`
	ReadTimeoutStr  string `yaml:"readTimeout"`
	WriteTimeoutStr string `yaml:"writeTimeout"`
	IdleTimeoutStr  string `yaml:"idleTimeout"`
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

type Route struct {
	Backend  string `yaml:"backend"`
	FilePath string `yaml:"filePath"`
}

type Proxy struct {
	Strategy      string `yaml:"strategy"`
	Matcher       string `yaml:"matcher"`
	NoRouteStatus int    `yaml:"noRoute"`
	MaxConn       int    `yaml:"maxConn"`
	BufferSize    int    `yaml:"bufSize"`

	ShutdownWaitStr          string `yaml:"shutdownWait"`
	DialTimeoutStr           string `yaml:"dialTimeout"`
	ResponseHeaderTimeoutStr string `yaml:"respHeaderTimeout"`
	KeepAliveTimeoutStr      string `yaml:"keepAliveTimeout"`
	FlushIntervalStr         string `yaml:"flushInterval"`

	ShutdownWait          time.Duration
	DialTimeout           time.Duration
	ResponseHeaderTimeout time.Duration
	KeepAliveTimeout      time.Duration
	FlushInterval         time.Duration
}

type Guard struct {
	Secret    string `yaml:"secret"`
	TimeRange int    `yaml:"timeRange"`
}

type Reactor struct {
	CircuitBreaker `yaml:"circuitBreaker"`
	RateLimiter    `yaml:"rateLimiter"`
	Replay         `yaml:"replay"`
}

type CircuitBreaker struct {
	ThresholdCount       int    `yaml:"thresholdCount"`
	ThresholdDurationStr string `yaml:"thresholdDuration"`
	BlockDurationStr     string `yaml:"blockDuration"`
}

type RateLimiter struct {
	ThresholdCount int    `yaml:"thresholdCount"`
	TTL            string `yaml:"ttl"`
	LockNum        int    `yaml:"lockNum"`
}

type Replay struct {
	TTL string `yaml:"ttl"`
}

type Recover struct {
	StackSize         int  `yaml:"stackSize"`
	DisableStackAll   bool `yaml:"disableStackAll"`
	DisablePrintStack bool `yaml:"disablePrintStack"`
}

type Consul struct {
	Address  string `yaml:"address"`
	Scheme   string `yaml:"scheme"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Etcd struct {
	Addresses  []string `yaml:addresses`
	ServiceDir string   `yaml:"serviceDir"`
}

type Admin struct {
	Listen string `yaml:"listen"`
}

type Metric struct {
	Listen string `yaml:"listen"`
}

type WebUI struct {
	Listen string `yaml:"listen"`
}

func DefaultConfig() Config {
	r := Route{
		Backend:  "file",
		FilePath: "route.cnf",
	}

	p := Proxy{
		Strategy:              "rnd",
		Matcher:               "prefix",
		NoRouteStatus:         404,
		MaxConn:               10000,
		BufferSize:            128,
		ShutdownWait:          5 * time.Second,
		DialTimeout:           30 * time.Second,
		FlushInterval:         time.Second,
		ResponseHeaderTimeout: time.Second,
		KeepAliveTimeout:      time.Second,
	}

	g := Guard{
		Secret:    "test123",
		TimeRange: 600,
	}

	reactor := Reactor{
		CircuitBreaker: CircuitBreaker{
			ThresholdCount:       60,
			ThresholdDurationStr: "1m",
			BlockDurationStr:     "1m",
		},
	}

	recover := Recover{
		StackSize: 4096,
	}

	listen := Listen{
		Address:      ":8088",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Second,
	}

	c := Config{
		Listen:  listen,
		Proxy:   p,
		Route:   r,
		Guard:   g,
		Reactor: reactor,
		Recover: recover,
	}

	return c
}

func Load(file string) (*Config, error) {
	var err error
	var c Config
	var f *os.File
	var b []byte
	f, err = os.Open(file)
	if err != nil {
		return nil, err
	}

	b, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	// parse string to duration
	c.Proxy.DialTimeout, err = time.ParseDuration(c.Proxy.DialTimeoutStr)
	if err != nil {
		return nil, err
	}
	c.Proxy.ShutdownWait, err = time.ParseDuration(c.Proxy.ShutdownWaitStr)
	if err != nil {
		return nil, err
	}
	c.Proxy.FlushInterval, err = time.ParseDuration(c.Proxy.FlushIntervalStr)
	if err != nil {
		return nil, err
	}
	c.Proxy.ResponseHeaderTimeout, err = time.ParseDuration(c.Proxy.ResponseHeaderTimeoutStr)
	if err != nil {
		return nil, err
	}
	c.Proxy.KeepAliveTimeout, err = time.ParseDuration(c.Proxy.KeepAliveTimeoutStr)
	if err != nil {
		return nil, err
	}

	c.Listen.IdleTimeout, err = time.ParseDuration(c.Listen.IdleTimeoutStr)
	if err != nil {
		return nil, err
	}
	c.Listen.ReadTimeout, err = time.ParseDuration(c.Listen.ReadTimeoutStr)
	if err != nil {
		return nil, err
	}
	c.Listen.WriteTimeout, err = time.ParseDuration(c.Listen.WriteTimeoutStr)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
