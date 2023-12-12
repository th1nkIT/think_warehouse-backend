package echokit

import (
	"fmt"

	"github.com/wit-id/blueprint-backend-go/toolkit/config"
)

func NewRuntimeConfig(cfg config.KVStore, path string) *RuntimeConfig {
	r := RuntimeConfig{}

	r.Port = cfg.GetInt(fmt.Sprintf("%s.port", path))
	r.RequestTimeoutConfig = &TimeoutConfig{
		Timeout: cfg.GetDuration(fmt.Sprintf("%s.request-timeout", path)),
	}
	r.ShutdownTimeoutDuration = cfg.GetDuration(fmt.Sprintf("%s.shutdown.timeout-duration", path))
	r.ShutdownWaitDuration = cfg.GetDuration(fmt.Sprintf("%s.shutdown.wait-duration", path))
	r.HealthCheckPath = cfg.GetString(fmt.Sprintf("%s.healthcheck-path", path))
	r.InfoCheckPath = cfg.GetString(fmt.Sprintf("%s.info-path", path))

	return &r
}
