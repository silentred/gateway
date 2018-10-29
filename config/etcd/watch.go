package etcd

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"encoding/json"

	etcd "github.com/coreos/etcd/client"
	"github.com/silentred/glog"
	"github.com/silentred/toolkit/util"
	"github.com/silentred/gateway/config"
	"github.com/silentred/gateway/guard"
	"github.com/silentred/gateway/reactor"
	"github.com/silentred/gateway/route"
)

// Backend is a service backend supported by etcd
type Backend struct {
	table  *route.Table
	cfg    *config.Etcd
	client etcd.Client
	kapi   etcd.KeysAPI
}

// NewBackend returns a new etcd backend
func NewBackend(t *route.Table, cfg *config.Etcd) *Backend {
	etcdCfg := etcd.Config{
		Endpoints:               cfg.Addresses,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	cli, err := etcd.New(etcdCfg)
	if err != nil {
		log.Fatalf("[etcd] init error: %v", err)
	}
	kapi := etcd.NewKeysAPI(cli)

	return &Backend{
		table:  t,
		cfg:    cfg,
		client: cli,
		kapi:   kapi,
	}
}

// Watch http service
func (b *Backend) Watch() {
	// list /iget/service/http, if KeyNotFound, create it
	dirs, err := b.listChildren(b.cfg.ServiceDir, true)
	if etcd.IsKeyNotFound(err) {
		err = b.createDir(b.cfg.ServiceDir)
		if err != nil {
			glog.Errorf("[etcd] create prefix dir [%s]: %v", b.cfg.ServiceDir, err)
			return
		}
	} else if err != nil {
		glog.Errorf("[etcd] prefix dir: %v", err)
		return
	}

	// if has sub dir, init them.
	// including: list each subdir, init services.
	if len(dirs) > 0 {
		// service
		for _, dirNode := range dirs {
			var targets []config.TargetInfo
			// target
			targets = b.fetchService(dirNode.Key)
			// make route, service, target
			glog.Debugf("[etcd] init services: %+v", targets)
			err = b.updateRouteService(targets)
			if err != nil {
				glog.Errorf("[etcd] err=%v", err)
			}
		}
	}

	// watch
	// found change. know which service
	// list service dir, get all targets, rebuild it
	for {
		resp, err := b.watchDir(b.cfg.ServiceDir)
		if err != nil {
			glog.Errorf("[etcd] watch err=%v; maybe etcd is offline", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var targets = b.fetchService(filepath.Dir(resp.Node.Key))
		glog.Debugf("[etcd] found services change: %+v", targets)
		err = b.updateRouteService(targets)
		if err != nil {
			glog.Errorf("[etcd] err=%v", err)
		}
	}
}

// SetTarget implements the Backend iface
func (b *Backend) SetTarget(t config.TargetInfo, h config.HealthCheck) error {
	var err error
	var jsonBytes []byte
	var targetAddr = fmt.Sprintf("%s:%d", t.TargetHost, t.TargetPort)
	var id = route.TargetID(t.ServiceName, targetAddr)
	var key = fmt.Sprintf("%s/%s/%s", b.cfg.ServiceDir, t.ServiceName, id)

	jsonBytes, err = json.Marshal(t)
	if err != nil {
		glog.Error(err)
		return err
	}

	_, err = b.kapi.Set(context.Background(), key, util.String(jsonBytes), &etcd.SetOptions{})
	if err != nil {
		glog.Errorf("[etcd] set key %s failed: %v", key, err)
	}
	return err
}

// DelTarget implements the Backend iface
func (b *Backend) DelTarget(t config.TargetInfo, id string) error {
	var err error
	var targetID = route.TargetID(t.ServiceName, t.TargetHost)
	var key = fmt.Sprintf("%s/%s/%s", b.cfg.ServiceDir, t.ServiceName, targetID)

	_, err = b.kapi.Delete(context.Background(), key, &etcd.DeleteOptions{})
	if err != nil {
		glog.Errorf("[etcd] del key %s failed: %v", key, err)
	}
	return err
}

func (b *Backend) fetchService(dir string) []config.TargetInfo {
	var targets []config.TargetInfo
	// target
	files, err := b.listChildren(dir, false)
	if err != nil {
		glog.Errorf("[etcd] ls dir [%s]: %v", dir, err)
		return targets
	}

	for _, file := range files {
		var target config.TargetInfo
		err = json.Unmarshal(util.Slice(file.Value), &target)
		if err != nil {
			glog.Errorf("[etcd] read target [%s]=%s err=%v", file.Key, file.Value, err)
			continue
		} else {
			targets = append(targets, target)
		}
	}

	return targets
}

func (b *Backend) updateRouteService(targetValues []config.TargetInfo) error {
	var r route.Route
	var s *route.Service
	var tmp config.TargetInfo
	var tmpTarget route.Target
	var targetID, host string

	// TODO: 需要容错
	for i := 0; i < len(targetValues); i++ {
		tmp = targetValues[i]
		// first
		host = fmt.Sprintf("%s:%d", tmp.TargetHost, tmp.TargetPort)
		targetID = getTargetID(tmp.ServiceName, host)
		tmpTarget = route.NewTarget(targetID, host, tmp.TargetWeight)
		if i == 0 {
			r = route.NewRoute(tmp.RouteHost, tmp.RoutePrefix)
			s = route.NewService(tmp.ServiceName, "http", tmp.ServiceStrip, tmpTarget, guard.DefaultGroup, reactor.DefaultGroup)
		} else {
			s.AddTarget(tmpTarget)
		}
	}

	if s == nil {
		return fmt.Errorf("[nil service] may be data error in etcd; targetValues: %+v", targetValues)
	}

	b.table.AddOrUpdate(r, s)
	return nil
}

func (b *Backend) watchDir(dir string) (*etcd.Response, error) {
	w := b.kapi.Watcher(dir, &etcd.WatcherOptions{
		Recursive: true,
	})

	return w.Next(context.Background())
}

func (b *Backend) listChildren(path string, dir bool) (etcd.Nodes, error) {
	resp, err := b.kapi.Get(context.Background(), path, &etcd.GetOptions{})
	if err != nil {
		return nil, err
	}

	var nodes = make(etcd.Nodes, 0, len(resp.Node.Nodes))
	for _, node := range resp.Node.Nodes {
		if (dir && node.Dir) || (!dir && !node.Dir) {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func (b *Backend) createDir(dir string) error {
	_, err := b.kapi.Set(context.Background(), dir, "", &etcd.SetOptions{
		Dir:       true,
		PrevExist: etcd.PrevNoExist,
	})

	return err
}

func (b *Backend) removeDir(dir string) error {
	_, err := b.kapi.Delete(context.Background(), dir, &etcd.DeleteOptions{
		Dir:       true,
		Recursive: true,
	})
	return err
}

// get hash id from serviceName and targetHost
func getTargetID(svcName, targetHost string) string {
	hash := sha1.Sum(util.Slice(svcName + targetHost))
	return fmt.Sprintf("%s-%x", svcName, hash[:10])
}
