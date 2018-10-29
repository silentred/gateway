package config

type Backend interface {
	Watch()
	SetTarget(t TargetInfo, h HealthCheck) error
	DelTarget(t TargetInfo, id string) error
}
