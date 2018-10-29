package consul

import (
	"log"
	"strconv"
	"strings"
	"time"

	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/silentred/glog"
	"github.com/silentred/gateway/config"
	"github.com/silentred/gateway/guard"
	"github.com/silentred/gateway/reactor"
	"github.com/silentred/gateway/route"
)

type Backend struct {
	cli          *api.Client
	table        *route.Table
	cfg          *config.Consul
	WaitDuration time.Duration
}

// NewBackend returns a new Backend based on Consul
func NewBackend(t *route.Table, cfg *config.Consul) *Backend {
	cli, err := NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &Backend{
		cli:          cli,
		table:        t,
		cfg:          cfg,
		WaitDuration: 2 * time.Second,
	}
}

// Watch implements the Backend iface
func (cb *Backend) Watch() {
	var q api.QueryOptions

	for {
		var routes []route.Route
		var services = make(map[string]*route.Service)
		var routeMap = make(map[uint32]string)

		var svcNames []string
		var checks api.HealthChecks
		var meta *api.QueryMeta
		var err error
		var svcEntries []*api.ServiceEntry

		var r route.Route
		var s *route.Service

		// for polling
		q.RequireConsistent = true
		checks, meta, err = cb.cli.Health().State("any", &q)
		if err != nil {
			glog.Errorf("[Consul] maybe offline: %v", err)
			time.Sleep(cb.WaitDuration)
			continue
		}
		q.WaitIndex = meta.LastIndex

		// get service name list
		for _, check := range checks {
			if !strings.Contains(check.CheckID, "serfHealth") && !inSliceString(svcNames, check.ServiceName) {
				svcNames = append(svcNames, check.ServiceName)
			}
		}

		// wait for service turning to passing status
		time.Sleep(cb.WaitDuration)

		for _, name := range svcNames {
			tmpQ := api.QueryOptions{}
			svcEntries, meta, err = cb.cli.Health().Service(name, "", true, &tmpQ)
			if err != nil {
				glog.Error(err)
				glog.Errorf("[Consul] maybe offline: %v", err)
				continue
			}
			q.WaitIndex = meta.LastIndex

			r, s = dealServiceEntries(svcEntries)

			if s != nil {
				// set to containers
				routes = append(routes, r)
				services[s.Name] = s
				routeMap[r.HashID] = s.Name
			}
		}

		for _, route := range routes {
			glog.Debugf("[Consul] route: %+v \n", route)
		}

		for _, item := range services {
			glog.Debugf("[Consul] svc: %+v targets:%+v\n", item, item.Targets.List())
		}

		for key, item := range routeMap {
			glog.Debugf("[Consul] route:%d -> svc:%s \n", key, item)
		}

		cb.table.SetAll(routes, services, routeMap)
		time.Sleep(cb.WaitDuration)
	}
}

// SetTarget implements the Backend interface
func (cb *Backend) SetTarget(t config.TargetInfo, h config.HealthCheck) error {
	var svcID string
	var tags []string
	var targetAddr = fmt.Sprintf("%s:%d", t.TargetHost, t.TargetPort)

	// TODO check params
	glog.Debugf("put service: %+v %+v", t, h)
	if t.ServiceName == "" || t.TargetHost == "" {
		return config.ErrBadParam
	}

	var tagMap = map[string]string{
		"gw.host":   t.RouteHost,
		"gw.prefix": t.RoutePrefix,
		"gw.strip":  t.ServiceStrip,
		"gw.weight": fmt.Sprintf("%d", t.TargetWeight),
	}

	for key, val := range tagMap {
		tags = append(tags, fmt.Sprintf("%s=%s", key, val))
	}

	svcID = route.TargetID(t.ServiceName, targetAddr)

	svcReg := &api.AgentServiceRegistration{
		ID:      svcID,
		Name:    t.ServiceName,
		Port:    t.TargetPort,
		Address: t.TargetHost,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			TCP:      h.HealthCheckAddr,
			Interval: h.HealthCheckInterval,
			Timeout:  "2s",
		},
	}

	return cb.cli.Agent().ServiceRegister(svcReg)
}

// DelTarget implements Backend interface
func (cb *Backend) DelTarget(t config.TargetInfo, id string) error {
	return cb.cli.Agent().ServiceDeregister(id)
}

func dealServiceEntries(entries []*api.ServiceEntry) (route.Route, *route.Service) {
	var r route.Route
	var s *route.Service
	var host, prefix string
	var strip, tHost, tWeight, tID string
	var tPort, tWeightInt int

	for i := 0; i < len(entries); i++ {
		var t route.Target
		// create target
		tWeight = tagValue(entries[i].Service.Tags, "gw.weight")
		tWeightInt, _ = strconv.Atoi(tWeight)

		// TODO make sure target is not empty
		tHost = entries[i].Service.Address
		tPort = entries[i].Service.Port
		tID = entries[i].Service.ID
		t = route.NewTarget(tID, fmt.Sprintf("%s:%d", tHost, tPort), tWeightInt)

		if i == 0 {
			// TODO make sure route and service.strip is not empty
			// create route
			host = tagValue(entries[i].Service.Tags, "gw.host")
			prefix = tagValue(entries[i].Service.Tags, "gw.prefix")
			r = route.NewRoute(host, prefix)

			// create service
			strip = tagValue(entries[i].Service.Tags, "gw.strip")
			s = route.NewService(entries[i].Service.Service, "http", strip, t, guard.DefaultGroup, reactor.DefaultGroup)
		} else {
			s.AddTarget(t)
		}
	}

	return r, s
}

func tagValue(tags []string, key string) string {
	for _, item := range tags {
		if strings.HasPrefix(item, key) {
			return item[len(key)+1:]
		}
	}
	return ""
}

func inSliceString(s []string, n string) bool {
	for i := 0; i < len(s); i++ {
		if n == s[i] {
			return true
		}
	}
	return false
}
