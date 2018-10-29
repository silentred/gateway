package route

import (
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/silentred/glog"
	"github.com/silentred/gateway/guard"
	"github.com/silentred/gateway/reactor"
)

const (
	// PickRoundRobin use round-robin to pick targets
	PickRoundRobin int = iota
	PickRandom
	PickWeight
)

// DirectorFunc change the request to target server
// type DirectorFunc func(*http.Request)

// Service represents backend service
type Service struct {
	Name     string         `json:"name"`
	Scheme   string         `json:"scheme"`
	Strip    string         `json:"strip"`
	Targets  TargetManager  `json:"targets"`
	Guards   *guard.Group   `json:"guards"`
	Reactors *reactor.Group `json:"reactors"`
	TotalCnt uint64         `json:"total_cnt"`
	mut      sync.Mutex
}

func NewService(name, scheme, strip string, t Target, guards *guard.Group, reactors *reactor.Group) *Service {
	if scheme == "" {
		scheme = "http"
	}
	ts := NewTargets(PickRoundRobin, t)
	s := &Service{
		Name:     name,
		Scheme:   scheme,
		Strip:    strip,
		Targets:  ts,
		Guards:   guards,
		Reactors: reactors,
		mut:      sync.Mutex{},
	}
	return s
}

func (s *Service) StripPrefix(path string) string {
	if strings.HasPrefix(path, s.Strip) {
		return path[len(s.Strip):]
	}
	return path
}

// Pick one target from targes of server
func (s *Service) Pick() (Target, error) {
	if len(s.Targets.List()) == 0 {
		return Target{}, ErrNoTarget
	}

	t := s.Targets.Pick(s.TotalCnt)
	return t, nil
}

// Director is able to modify a http request for ReverseProxy
func (s *Service) Director(req *http.Request) {
	t, err := s.Pick()
	glog.Debugf("[Svc] picked target=%s w=%d err=%v", t.Host, t.Weight, err)
	if err != nil {
		glog.Errorf("[Svc] picked target=%s w=%d err=%v", t.Host, t.Weight, err)
		return
	}
	// increase total
	atomic.AddUint64(&s.TotalCnt, 1)

	req.URL.Scheme = s.Scheme
	req.URL.Path = s.StripPrefix(req.URL.Path)
	req.URL.Host = t.Host
	req.Host = t.Host

	if _, ok := req.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
		req.Header.Set("User-Agent", "entree/v0.0.2")
	}
}

// Modify is a ModifyResponse func, which could observe and modify http response.
func (s *Service) Modify(resp *http.Response) error {
	var err error
	if s.Reactors != nil {
		err = s.Reactors.Modify(resp)
	}
	return err
}

// ObserveError implements Reator interface
func (s *Service) ObserveError(req *http.Request, resp *http.Response, err error) {
	if s.Reactors != nil {
		s.Reactors.ObserveError(req, resp, err)
	}
}

// AddTarget adds target to Targets
func (s *Service) AddTarget(t Target) {
	s.Targets.Add(t)
}

// DelTarget removes target from Targets
func (s *Service) DelTarget(host string) {
	s.Targets.Del(host)
}

// // MakeDirectorFunc returns a DirectorFunc
// func (s *Service) MakeDirectorFunc() DirectorFunc {
// 	t, err := s.Pick()
// 	if err != nil {
// 		// TODO log
// 		return nil
// 	}
// 	var scheme = s.Scheme
// 	var stripStr = s.Strip

// 	return func(req *http.Request) {
// 		req.URL.Scheme = scheme
// 		req.URL.Path = strip(req.URL.Path, stripStr)
// 		req.URL.Host = t.Host
// 		if _, ok := req.Header["User-Agent"]; !ok {
// 			// explicitly disable User-Agent so it's not set to default value
// 			req.Header.Set("User-Agent", "entree")
// 		}
// 	}
// }

// func strip(path, strip string) string {
// 	if strings.HasPrefix(path, strip) {
// 		return path[len(strip):]
// 	}
// 	return path
// }
