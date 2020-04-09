package admin

import (
	"encoding/json"
	"testing"

	"github.com/silentred/gateway/route"
	"github.com/stretchr/testify/assert"
)

func TestAdmin(t *testing.T) {
	// dto := routesDTO{
	// 	RouteHost:           "svc.baidu.com",
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
			Host:   "svc.baidu.com",
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
