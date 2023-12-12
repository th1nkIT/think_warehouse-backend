package postgres

import (
	"database/sql"
	"fmt"

	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/db"
)

func NewFromConfig(cfg config.KVStore, path string) (*sql.DB, error) {
	connOpt := db.DefaultConnectionOption()

	if maxIdle := cfg.GetInt(fmt.Sprintf("%s.conn.max-idle", path)); maxIdle > 0 {
		connOpt.MaxIdle = cfg.GetInt(fmt.Sprintf("%s.conn.max-idle", path))
	}

	if maxOpen := cfg.GetInt(fmt.Sprintf("%s.conn.max-open", path)); maxOpen > 0 {
		connOpt.MaxOpen = maxOpen
	}

	if maxLifetime := cfg.GetDuration(fmt.Sprintf("%s.conn.max-lifetime", path)); maxLifetime > 0 {
		connOpt.MaxLifetime = maxLifetime
	}

	if connTimeout := cfg.GetDuration(fmt.Sprintf("%s.conn.timeout", path)); connTimeout > 0 {
		connOpt.ConnectTimeout = connTimeout
	}

	if keepAlive := cfg.GetDuration(fmt.Sprintf("%s.conn.keep-alive-interval", path)); keepAlive > 0 {
		connOpt.KeepAliveCheckInterval = keepAlive
	}

	opt, err := db.NewDatabaseOption(
		cfg.GetString(fmt.Sprintf("%s.host", path)),
		cfg.GetInt(fmt.Sprintf("%s.port", path)),
		cfg.GetString(fmt.Sprintf("%s.username", path)),
		cfg.GetString(fmt.Sprintf("%s.password", path)),
		cfg.GetString(fmt.Sprintf("%s.schema", path)),
		connOpt,
	)
	if err != nil {
		return nil, err
	}

	return NewPostgresDatabase(opt)
}
