package config

import "github.com/silentred/gateway/util"

var (
	ErrBadParam = util.NewError(110110, "bad param of setting target")
)

// TargetInfo is the full info for create a route-service pair
type TargetInfo struct {
	RouteHost   string `json:"route_host"`
	RoutePrefix string `json:"route_prefix"`

	ServiceName  string `json:"service_name"`
	ServiceStrip string `json:"service_strip"`

	TargetHost   string `json:"target_host"`
	TargetPort   int    `json:"target_port"`
	TargetWeight int    `json:"target_weight"`
}

// HealthCheck is for consul
type HealthCheck struct {
	HealthCheckAddr     string `json:"hc_addr"`
	HealthCheckInterval string `json:"hc_interval"`
}
