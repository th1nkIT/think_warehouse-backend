package redis

import (
	"fmt"

	"think_warehouse/toolkit/config"
	"think_warehouse/toolkit/db"

	"github.com/go-redis/redis/v8"
)

func NewFromConfig(cfg config.KVStore, path string) (*redis.Client, error) {
	opt, err := db.NewRedisOption(
		cfg.GetString(fmt.Sprintf("%s.host", path)),
		cfg.GetInt(fmt.Sprintf("%s.port", path)),
		cfg.GetString(fmt.Sprintf("%s.password", path)),
	)
	if err != nil {
		return nil, err
	}

	return NewRedisDatabase(opt)
}
