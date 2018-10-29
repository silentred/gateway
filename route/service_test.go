package route

import (
	"testing"

	"github.com/silentred/gateway/guard"
	"github.com/silentred/gateway/reactor"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	var err error
	target := NewTarget("", "localhost:8080", 1)
	service := NewService("test service", "", "/v1", target, nil, nil)
	assert.NotNil(t, service.Targets)

	pickedTarget, err := service.Pick()
	assert.NoError(t, err)
	assert.Equal(t, target.Host, pickedTarget.Host)

	path := service.StripPrefix("/v1/hello")
	assert.Equal(t, "/hello", path)

	newTarget := NewTarget("", "localhost:8081", 2)
	service.AddTarget(newTarget)
	assert.Equal(t, 2, len(service.Targets.List()))

	service.Guards = guard.NewGroup("gGroup", &guard.MockGuard{}, &guard.MockGuard{})

	service.Reactors = reactor.NewGroup("rGroup", &reactor.MockReactor{})

	// test json
	b, err := json.Marshal(service)
	assert.NoError(t, err)
	t.Logf("json: %s", string(b))

	service.DelTarget("localhost:8081")
	assert.Equal(t, 1, len(service.Targets.List()))
}
