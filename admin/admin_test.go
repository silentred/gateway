package admin

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/silentred/gateway/route"
)

func TestAdmin(t *testing.T) {
	// dto := routesDTO{
	// 	RouteHost:           "svc.luojilab.com",
	// 	RoutePrefix:         "/v1/hello",
	// 	ServiceName:         "testSvc",
	// 	ServiceStrip:        "/v1",
	// 	TargetHost:          "localhost",
	// 	TargetPort:          8080,
	// 	TargetWeight:        1,
	// 	HealthCheckAddr:     "localhost:8080",
	// 	HealthCheckInterval: "5s",
	// }

	dto := routesDTO{
		Route: route.Route{
			Host:   "svc.luojilab.com",
			Prefix: "/v1/hello",
		},
		Service: &route.Service{
			Name:  "testSvc",
			Strip: "/v1",
		},
	}

	body, err := json.Marshal(dto)
	assert.NoError(t, err)
	t.Log(string(body))

}
