package main

import (
	"time"

	"think_warehouse/common/echohttp"
	"think_warehouse/common/httpservice"
	"think_warehouse/toolkit/db/postgres"
	"think_warehouse/toolkit/log"
	"think_warehouse/toolkit/runtimekit"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func main() {
	var err error

	setDefaultTimezone()

	appContext, cancel := runtimekit.NewRuntimeContext()
	defer func() {
		cancel()

		if err != nil {
			log.FromCtx(appContext).Error(err, "found error")
		}
	}()

	// Set config file (env)
	appConfig, err := envConfigVariable("config.yaml")
	if err != nil {
		return
	}

	// setup db
	mainDB, err := postgres.NewFromConfig(appConfig, "db")
	if err != nil {
		return
	}

	// setup redis db  (DEFAULT DISABLE)
	// redisDB, err := redis.NewFromConfig(appConfig, "redis")
	// if err != nil {
	//	 return
	// }

	// setup logging
	logger, err := log.NewFromConfig(appConfig, "log")
	if err != nil {
		return
	}

	logger.Set()

	// setup service
	svc := httpservice.NewService(mainDB, appConfig)

	// expose echo http server
	echohttp.RunEchoHTTPService(appContext, svc, appConfig)
}

func setDefaultTimezone() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		loc = time.Now().Location()
	}

	time.Local = loc
}

func envConfigVariable(filePath string) (cfg *viper.Viper, err error) {
	cfg = viper.New()
	cfg.SetConfigFile(filePath)

	if err = cfg.ReadInConfig(); err != nil {
		err = errors.Wrap(err, "Error while reading config file")

		return
	}

	return
}
