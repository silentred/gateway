package config

import "testing"
import "encoding/json"
import "github.com/stretchr/testify/assert"

func TestJson(t *testing.T) {
	var target TargetInfo
	var health HealthCheck

	var body = []byte(`{"route_host": "test.luoji.com", "hc_addr": "127.0.0.1:8080"}`)

	err := json.Unmarshal(body, &struct {
		*TargetInfo
		*HealthCheck
	}{&target, &health})

	assert.NoError(t, err)
	assert.Equal(t, "test.luoji.com", target.RouteHost)
	assert.Equal(t, "127.0.0.1:8080", health.HealthCheckAddr)

}
