package httpservice

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/wit-id/blueprint-backend-go/common/constants"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

// RoutesPayload ...
type RoutesPayload struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Name   string `json:"name"`
}

// RouteAccess ...
type RouteAccess struct {
	Access string `json:"access"`
	Path   string `json:"path"`
}

// SetRouteConfig ...
func SetRouteConfig(ctx context.Context, svc *Service, cfg config.KVStore, e *echo.Echo) {
	s := sqlc.New(svc.mainDB)

	route, _ := json.MarshalIndent(e.Routes(), "", " ")

	var (
		routesMap    []RoutesPayload
		routesAccess []RouteAccess
	)

	err := json.Unmarshal(route, &routesMap)
	if err != nil {
		log.FromCtx(ctx).Info("failed unmarshal routes map")
	}

	for v := range routesMap {
		if routesMap[v].Method == "HEAD" {
			routesAccess = append(routesAccess, RouteAccess{
				Access: constants.AccessView + "|" + constants.AccessCreate + "|" + constants.AccessUpdate + "|" + constants.AccessGet + "|" + constants.AccessList + "|" + constants.AccessDelete,
				Path:   routesMap[v].Path,
			})
		}
	}

	routeRaw, err := json.Marshal(routesAccess)
	if err != nil {
		log.FromCtx(ctx).Info("failed marshal route access json")
	}

	_, err = s.SetConfig(ctx, sqlc.SetConfigParams{
		Key: cfg.GetString(constants.ConfigRoutes),
		Description: sql.NullString{
			String: "config routes for access",
			Valid:  true,
		},
		Value: string(routeRaw),
	})

	if err != nil {
		log.Println("failed set config routes access")
	} else {
		log.Println("successfully set config routes access")
	}
}
