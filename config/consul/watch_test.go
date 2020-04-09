package consul

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/silentred/gateway/config"
	"github.com/silentred/gateway/route"
	"github.com/stretchr/testify/assert"
)

func TestListService(t *testing.T) {
	// cli, err := api.NewClient(api.DefaultConfig())
	// assert.NoError(t, err)

	// list all service
	// list, err := cli.Agent().Services()
	// assert.NoError(t, err)
	// for key, val := range list {
	// 	t.Logf("key:%s val:%#v \n", key, val)
	// }

	// list health
	// q := api.QueryOptions{}
	// svcEntries, meta, err := cli.Health().Service("svcA", "", true, &q)
	// assert.NoError(t, err)
	// t.Logf("meta:%#v \n", meta)
	// for _, item := range svcEntries {
	// 	t.Logf("svc_entry:%#v \n", item.Service)
	// }

	// q := api.QueryOptions{
	// 	RequireConsistent: true,
	// 	//WaitIndex:         0x78,
	// }
	// check, meta, err := cli.Health().State("any", &q)
	// assert.NoError(t, err)
	// t.Logf("meta:%#v \n", meta)
	// for _, item := range check {
	// 	t.Logf("check:%#v \n", item)
	// }

	// catalog
	// svcs, meta, err := cli.Catalog().Service("svcA", "", &q)
	// assert.NoError(t, err)
	// t.Logf("meta:%#v \n", meta)
	// for _, item := range svcs {
	// 	t.Logf("svc:%#v \n", item)
	// }

}

func testAddService(t *testing.T) {
	cli, err := api.NewClient(api.DefaultConfig())
	assert.NoError(t, err)
	addSvc(t, cli)
}

func addSvc(t *testing.T, cli *api.Client) {
	svc := &api.AgentServiceRegistration{
		ID:      "svcC02",
		Name:    "svcC",
		Port:    8080,
		Address: "localhost",
		Tags: []string{
			"gw.host=www.baidu.com",
			"gw.prefix=/v1/C/hello",
			"gw.strip=/v1",
			"gw.weight=20",
		},
		Check: &api.AgentServiceCheck{
			TCP:      "localhost:8080",
			Interval: "5s",
			Timeout:  "1s",
		},
	}

	err := cli.Agent().ServiceRegister(svc)
	assert.NoError(t, err)
}

func delSvc(t *testing.T, cli *api.Client) {
	cli.Agent().ServiceDeregister("svcA01")
}

func testKV(t *testing.T) {
	cli, err := api.NewClient(api.DefaultConfig())
	assert.NoError(t, err)

	var lastIndex uint64
	for {
		q := api.QueryOptions{
			WaitIndex: lastIndex,
		}
		kvs, meta, err := cli.KV().List("gw.luoji", &q)
		assert.NoError(t, err)
		fmt.Printf("meta : %+v, kvs: %+v \n", meta, kvs)
		for _, item := range kvs {
			fmt.Printf("kv: %+v \n", item)
		}
		lastIndex = meta.LastIndex

		time.Sleep(time.Second / 2)
		fmt.Println("round")
	}

}

// real test of watching services in consul
func testConsulBackend(t *testing.T) {
	//cli, err := api.NewClient(api.DefaultConfig())
	//assert.NoError(t, err)

	table := route.NewTable()
	config := &config.Consul{}
	cb := NewBackend(table, config)

	cb.Watch()
}
