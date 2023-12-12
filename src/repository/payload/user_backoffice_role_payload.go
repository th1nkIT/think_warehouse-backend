package payload

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/constants"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"
)

type UserBackofficeRolePayload struct {
	Name        string                      `json:"name" valid:"required"`
	Access      []ConfigRouteAccessResponse `json:"access" valid:"required"`
	IsAllAccess bool                        `json:"is_all_access"`
}

type ListUserBackofficeRolePayload struct {
	Filter ListUserBackofficeRoleFilterPayload `json:"filter"`
	Limit  int32                               `json:"limit" valid:"required"`
	Offset int32                               `json:"page" valid:"required"`
	Order  string                              `json:"order" valid:"required"`
	Sort   string                              `json:"sort" valid:"required"` // ASC, DESC
}

type ListUserBackofficeRoleFilterPayload struct {
	SetName bool   `json:"set_name"`
	Name    string `json:"name"`
}

type readUserBackofficeRoleDataPayload struct {
	ID          int64                       `json:"id"`
	Name        string                      `json:"name"`
	Access      []ConfigRouteAccessResponse `json:"access"`
	IsAllAccess bool                        `json:"is_all_access"`
	CreatedAt   time.Time                   `json:"created_at"`
	CreatedBy   string                      `json:"created_by"`
	UpdatedAt   *time.Time                  `json:"updated_at"`
	UpdatedBy   *string                     `json:"updated_by"`
}

func (payload *UserBackofficeRolePayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *ListUserBackofficeRolePayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *UserBackofficeRolePayload) ToEntityCreate(userData sqlc.GetUserBackofficeRow) (data sqlc.InsertUserBackofficeRoleParams) {
	data = sqlc.InsertUserBackofficeRoleParams{
		Name: payload.Name,
		IsAllAccess: sql.NullBool{
			Bool:  payload.IsAllAccess,
			Valid: true,
		},
		CreatedBy: userData.Guid,
	}

	if !payload.IsAllAccess {
		access, _ := json.Marshal(payload.Access)

		data.Access = sql.NullString{
			String: string(access),
			Valid:  true,
		}
	}

	return
}

func (payload *UserBackofficeRolePayload) ToEntityUpdate(id int64, userData sqlc.GetUserBackofficeRow) (data sqlc.UpdateUserBackofficeRoleParams) {
	data = sqlc.UpdateUserBackofficeRoleParams{
		Name: payload.Name,
		IsAllAccess: sql.NullBool{
			Bool:  payload.IsAllAccess,
			Valid: true,
		},
		UpdatedBy: sql.NullString{
			String: userData.Guid,
			Valid:  true,
		},
		ID: id,
	}

	if !payload.IsAllAccess {
		access, _ := json.Marshal(payload.Access)

		data.Access = sql.NullString{
			String: string(access),
			Valid:  true,
		}
	}

	return
}

func (payload *ListUserBackofficeRolePayload) ToEntity() (data sqlc.ListUserBackofficeRoleParams) {
	orderParam := constants.DefaultOrderValue

	data = sqlc.ListUserBackofficeRoleParams{
		SetName:   payload.Filter.SetName,
		Name:      "%" + payload.Filter.Name + "%",
		LimitData: payload.Limit,
	}

	if payload.Limit == 0 {
		data.LimitData = 10
	}

	if payload.Offset == 0 {
		data.OffsetPage = (1 * data.LimitData) - data.LimitData
	} else {
		data.OffsetPage = (payload.Offset * data.LimitData) - data.LimitData
	}

	if payload.Order != "" {
		orderParam = payload.Order + " ASC"

		if payload.Sort != "" {
			orderParam = payload.Order + " " + payload.Sort
		}
	}

	data.OrderParam = orderParam

	return
}

func ToPayloadUserBackofficeRole(userBackofficeRole sqlc.UserBackofficeRole) (payload readUserBackofficeRoleDataPayload) {
	payload = readUserBackofficeRoleDataPayload{
		ID:          userBackofficeRole.ID,
		Name:        userBackofficeRole.Name,
		IsAllAccess: userBackofficeRole.IsAllAccess.Bool,
		CreatedAt:   userBackofficeRole.CreatedAt,
		CreatedBy:   userBackofficeRole.CreatedBy,
		UpdatedAt:   nil,
		UpdatedBy:   nil,
	}

	if userBackofficeRole.Access.Valid {
		var access []ConfigRouteAccessResponse

		err := json.Unmarshal([]byte(userBackofficeRole.Access.String), &access)
		if err != nil {
			log.FromCtx(context.Background()).Info("failed unmarshal data route access")
		}

		payload.Access = access
	}

	if userBackofficeRole.UpdatedAt.Valid {
		updatedAt := userBackofficeRole.UpdatedAt.Time
		payload.UpdatedAt = &updatedAt
	}

	if userBackofficeRole.UpdatedBy.Valid {
		updatedBy := userBackofficeRole.UpdatedBy.String
		payload.UpdatedBy = &updatedBy
	}

	return
}

func ToPayloadListUserBackofficeRole(userBackofficeRole []sqlc.UserBackofficeRole) (payload []*readUserBackofficeRoleDataPayload) {
	payload = make([]*readUserBackofficeRoleDataPayload, len(userBackofficeRole))

	for i := range userBackofficeRole {
		payload[i] = new(readUserBackofficeRoleDataPayload)
		data := ToPayloadUserBackofficeRole(userBackofficeRole[i])
		payload[i] = &data
	}

	return
}
