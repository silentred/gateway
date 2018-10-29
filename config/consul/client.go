package consul

import (
	"log"

	"github.com/hashicorp/consul/api"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/silentred/gateway/config"
)

func NewClient(c *config.Consul) (*api.Client, error) {
	var scheme string
	var basicAuth *api.HttpBasicAuth

	scheme = c.Scheme
	if scheme == "" {
		scheme = "http"
	}

	if c.Address == "" {
		log.Fatal("consul address cannot be empty")
	}

	if c.Username != "" {
		basicAuth = &api.HttpBasicAuth{
			Username: c.Username,
			Password: c.Password,
		}
	}

	cfg := api.Config{
		Address:   c.Address,
		Scheme:    scheme,
		Transport: cleanhttp.DefaultPooledTransport(),
		HttpAuth:  basicAuth,
	}

	return api.NewClient(&cfg)
}
