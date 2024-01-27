package service_test

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"think_warehouse/src/auth_token/service"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"think_warehouse/toolkit/config"
)

func TestAuthTokenService_ReadAuthToken(t *testing.T) {
	type fields struct {
		mainDB *sql.DB
		cfg    config.KVStore
	}

	type args struct {
		ctx     context.Context
		request sqlc.GetAuthTokenParams
	}

	tests := []struct {
		name          string
		fields        fields
		args          args
		wantAuthToken sqlc.AuthToken
		wantErr       bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewAuthTokenService(tt.fields.mainDB, tt.fields.cfg)
			gotAuthToken, err := s.ReadAuthToken(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAuthToken, tt.wantAuthToken) {
				t.Errorf("ReadAuthToken() gotAuthToken = %v, want %v", gotAuthToken, tt.wantAuthToken)
			}
		})
	}
}
