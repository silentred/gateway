package route

import (
	"net/http"
	"strings"
	"sync"
)

// Table containers the relations between http request and backend service
type Table struct {
	routes   []Route
	routeMap map[uint32]string
	services map[string]*Service
	root     *node
	mut      sync.RWMutex
}

// NewTable returns a route table
func NewTable() *Table {
	t := &Table{
		routes:   make([]Route, 0, 100),
		routeMap: make(map[uint32]string),
		services: make(map[string]*Service),
		root:     newNode(""),
		mut:      sync.RWMutex{},
	}
	return t
}

// Reset route table
func (rt *Table) Reset() error {
	rt.mut.Lock()
	defer rt.mut.Unlock()
	rt.routes = rt.routes[0:0]
	rt.routeMap = make(map[uint32]string)
	rt.services = make(map[string]*Service)
	return nil
}

// Add route and service to route table
func (rt *Table) Add(r Route, s *Service) error {
	rt.mut.Lock()
	defer rt.mut.Unlock()

	var id = HashID(r.Host, r.Prefix)
	if _, has := rt.routeMap[id]; has {
		return ErrDupRoute
	}

	rt.routes = append(rt.routes, r)
	rt.routeMap[id] = s.Name

	if _, has := rt.services[s.Name]; !has {
		rt.services[s.Name] = s
	}
	// add to trie
	rt.root.insert(r.Host + getCanonicalPrefix(r.Prefix))

	return nil
}

// AddOrUpdate route and service
func (rt *Table) AddOrUpdate(r Route, s *Service) {
	rt.mut.Lock()
	defer rt.mut.Unlock()

	var id = HashID(r.Host, r.Prefix)
	// add
	if _, has := rt.routeMap[id]; !has {
		rt.routes = append(rt.routes, r)
		rt.root.insert(r.Host + getCanonicalPrefix(r.Prefix))
	}
	// update
	rt.routeMap[id] = s.Name
	rt.services[s.Name] = s
}

// Del route from route table
func (rt *Table) Del(r *Route) bool {
	rt.mut.Lock()
	defer rt.mut.Unlock()

	var id = HashID(r.Host, r.Prefix)
	var cnt int

	if _, has := rt.routeMap[id]; !has {
		return false
	}

	// del from slice
	for idx, item := range rt.routes {
		if item.Host == r.Host && item.Prefix == r.Prefix {
			rt.routes = rt.routes[:idx+copy(rt.routes[idx:], rt.routes[idx+1:])]
		}
	}

	// try to remove from services
	svcToDel := rt.routeMap[id]
	for key := range rt.services {
		if key == svcToDel {
			cnt++
		}
	}
	if cnt <= 1 {
		delete(rt.services, svcToDel)
	}

	// remove from routeMap
	delete(rt.routeMap, id)

	return true
}

// Find service from route table
func (rt *Table) Find(r Route) *Service {
	rt.mut.RLock()
	defer rt.mut.RUnlock()

	var id = HashID(r.Host, r.Prefix)
	if name, has := rt.routeMap[id]; has {
		if svc, has := rt.services[name]; has {
			return svc
		}
	}
	return nil
}

// find service from route table
func (rt *Table) findByID(hashID uint32) *Service {
	rt.mut.RLock()
	defer rt.mut.RUnlock()

	if name, has := rt.routeMap[hashID]; has {
		if svc, has := rt.services[name]; has {
			return svc
		}
	}
	return nil
}

// FindByRequest do find service according to http request
func (rt *Table) FindByRequest(r *http.Request) *Service {
	rt.mut.RLock()
	defer rt.mut.RUnlock()

	var host = r.Host
	var uri string
	var id uint32

	if host == "" {
		host = r.Header.Get("Host")
	}

	// use trie tree to find prefix
	uri = rt.root.find(host + r.URL.Path)
	id = HashCrc32(uri)
	return rt.findByID(id)

	// use map to find prefix, not good
	// for _, item := range rt.routes {
	// 	//fmt.Printf("host:%s path:%s route-host:%s route-Prefix:%s \n", host, r.URL.Path, item.Host, item.Prefix)
	// 	if item.Host == host && strings.HasPrefix(r.URL.Path, item.Prefix) {
	// 		return rt.Find(item)
	// 	}
	// }
	// return nil
}

func (rt *Table) Services() map[string]*Service {
	return rt.services
}

func (rt *Table) SetAll(routes []Route, services map[string]*Service, routeMap map[uint32]string) {
	rt.mut.Lock()
	defer rt.mut.Unlock()

	rt.routes = routes
	rt.services = services
	rt.routeMap = routeMap

	for i := 0; i < len(routes); i++ {
		rt.root.insert(routes[i].Host + getCanonicalPrefix(routes[i].Prefix))
	}
}

func (rt *Table) Routes() map[Route]*Service {
	var routes = make(map[Route]*Service)

	for _, route := range rt.routes {
		if name, has := rt.routeMap[route.HashID]; has {
			if svc, hasSvc := rt.services[name]; hasSvc {
				routes[route] = svc
			}
		}
	}

	return routes
}

func getCanonicalPrefix(prefix string) string {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	return prefix
}
