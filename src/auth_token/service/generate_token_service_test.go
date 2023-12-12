package service_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/wit-id/blueprint-backend-go/src/auth_token/service"

	"github.com/wit-id/blueprint-backend-go/src/repository/payload"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
)

func TestAuthTokenService_AuthToken(t *testing.T) {
	tests := []struct {
		name          string
		request       payload.AuthTokenPayload
		wantAuthToken sqlc.AuthToken
		wantErr       bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			kvStore := viper.New()
			kvStore.Set("jwt.key", "token-key")
			kvStore.Set("jwt.expired", "24h")
			kvStore.Set("jwt.refresh_expired", "168h")

			s := service.NewAuthTokenService(nil, kvStore)
			gotAuthToken, err := s.AuthToken(ctx, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAuthToken, tt.wantAuthToken) {
				t.Errorf("AuthToken() gotAuthToken = %v, want %v", gotAuthToken, tt.wantAuthToken)
			}
		})
	}
}
