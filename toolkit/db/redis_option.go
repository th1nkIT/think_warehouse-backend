package db

import "github.com/pkg/errors"

type RedisOption struct {
	Host     string
	Port     int
	Password string
}

func NewRedisOption(host string, port int, psswd string) (*RedisOption, error) {
	if host == "" || port == 0 {
		return nil, errors.Wrapf(errInvalidDBSource, "db: host=%s port=%d", host, port)
	}

	return &RedisOption{
		Host:     host,
		Port:     port,
		Password: psswd,
	}, nil
}
