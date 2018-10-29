package etcd

import (
	"testing"

	etcd "github.com/coreos/etcd/client"
	"github.com/stretchr/testify/assert"
	"github.com/silentred/gateway/config"
)

func testEtcd(t *testing.T) {
	cfg := &config.Etcd{
		Addresses:  []string{"http://127.0.0.1:2379"},
		ServiceDir: "/iget/service/http",
	}
	back := NewBackend(nil, cfg)

	nodes, err := back.listChildren("iget/service/http", true)
	assert.Error(t, err)
	t.Logf("nodes : %v, err : %#v", nodes, err)
	assert.True(t, etcd.IsKeyNotFound(err))
	if e, ok := err.(etcd.Error); ok {
		assert.Equal(t, etcd.ErrorCodeKeyNotFound, e.Code)
	}

	err = back.createDir("iget/service/http")
	assert.NoError(t, err)

	err = back.createDir("iget/service/http/hello")
	assert.NoError(t, err)

	// watch
	resp, err := back.watchDir("iget/service/http")
	if err != nil {
		t.Logf("err:%v", err)
	}
	t.Logf("resp.Node: %v", resp.Node)
	// Node 可以得到 key, 判断是哪个service的变化

	err = back.removeDir("iget/service/http")
	assert.NoError(t, err)
}
