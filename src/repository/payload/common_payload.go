package payload

import (
	"context"
	"encoding/json"
	"strings"

	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

// ConfigRouteAccessResponse ...
type ConfigRouteAccessResponse struct {
	Page    string   `json:"page"`
	KeyPage string   `json:"key_page"`
	Access  []string `json:"access"`
}

// Images ...
type Images struct {
	Thumbnail string   `json:"thumbnail"`
	Image     []string `json:"image"`
}

// ConfigRoute ...
type ConfigRoute struct {
	Access string `json:"access"`
	Path   string `json:"path"`
}

// ToPayloadConfigRouteAccess ...
func ToPayloadConfigRouteAccess(prefixPath string, data sqlc.Config) (response []ConfigRouteAccessResponse) {
	var routes []ConfigRoute

	err := json.Unmarshal([]byte(data.Value), &routes)
	if err != nil {
		log.FromCtx(context.Background()).Info("failed unmarshal response value route access")
	}

	for v := range routes {
		if strings.Contains(routes[v].Path, prefixPath) && !strings.Contains(routes[v].Path, "*") {
			pageSplit := strings.Split(strings.ReplaceAll(routes[v].Path, prefixPath, ""), "/")
			pagesName := strings.ReplaceAll(strings.Join(pageSplit, " "), "-", "")

			response = append(response, ConfigRouteAccessResponse{
				Page:    strings.ToUpper(pagesName[:1]) + pagesName[1:],
				KeyPage: routes[v].Path,
				Access:  strings.Split(routes[v].Access, "|"),
			})
		}
	}

	return
}
