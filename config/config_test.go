package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	file := "../config.yaml"
	c, err := Load(file)
	assert.NoError(t, err)
	assert.Equal(t, ":8088", c.Listen.Address)
	assert.Equal(t, "rnd", c.Proxy.Strategy)
	assert.Equal(t, 60, c.Reactor.CircuitBreaker.ThresholdCount)
	assert.Equal(t, 4096, c.Recover.StackSize)
	assert.Equal(t, ":9088", c.Admin.Listen)
	assert.Equal(t, ":7088", c.Metric.Listen)
	assert.Equal(t, ":8090", c.WebUI.Listen)
	assert.Equal(t, "127.0.0.1:8500", c.Consul.Address)
	assert.True(t, len(c.Etcd.Addresses) > 0)
	t.Log(c)

	_ = DefaultConfig()
}
