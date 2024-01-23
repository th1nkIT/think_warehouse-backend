package payload

import (
	"database/sql"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/constants"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/common/utility"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"time"
)

type RegisterWarehousePayload struct {
	WarehouseCode string `json:"warehouse_code" valid:"required"`
	Name          string `json:"name" valid:"required"`
	Address       string `json:"address" valid:"required"`
	PhoneNumber   string `json:"phone_number" valid:"required"`
}

type UpdateWarehousePayload struct {
	Name        string `json:"name" valid:"required"`
	Address     string `json:"address" valid:"required"`
	PhoneNumber string `json:"phone_number" valid:"required"`
}

type ListWarehousePayload struct {
	Filter ListWarehouseFilterPayload `json:"filter"`
	Limit  int32                      `json:"limit" valid:"required"`
	Offset int32                      `json:"page" valid:"required"`
	Order  string                     `json:"order" valid:"required"`
	Sort   string                     `json:"sort" valid:"required"` // ASC, DESC
}

type ListWarehouseFilterPayload struct {
	SetName          bool   `json:"set_name"`
	Name             string `json:"name"`
	SetWarehouseCode bool   `json:"set_warehouse_code"`
	WarehouseCode    string `json:"warehouse_code"`
	SetActive        bool   `json:"set_active"`
	Active           string `json:"active"`
}

type readRegisterWarehousePayload struct {
	GUID          string                   `json:"id"`
	WarehouseCode string                   `json:"warehouse_code" valid:"required"`
	Name          string                   `json:"name" valid:"required"`
	Address       string                   `json:"address" valid:"required"`
	PhoneNumber   string                   `json:"phone_number" valid:"required"`
	Status        string                   `json:"status"`
	CreatedAt     time.Time                `json:"created_at"`
	CreatedBy     readUserWarehousePayload `json:"created_by"`
}

type readUserWarehousePayload struct {
	GUID string `json:"id"`
	Name string `json:"name"`
}

type readUpdateWarehousePayload struct {
	GUID          string                   `json:"id"`
	WarehouseCode string                   `json:"warehouse_code" valid:"required"`
	Name          string                   `json:"name" valid:"required"`
	Address       string                   `json:"address" valid:"required"`
	PhoneNumber   string                   `json:"phone_number" valid:"required"`
	Status        string                   `json:"status"`
	UpdatedAt     time.Time                `json:"updated_at"`
	UpdatedBy     readUserWarehousePayload `json:"updated_by"`
}

type readWarehousePayload struct {
	GUID          string                    `json:"id"`
	WarehouseCode string                    `json:"warehouse_code" valid:"required"`
	Name          string                    `json:"name" valid:"required"`
	Address       string                    `json:"address" valid:"required"`
	PhoneNumber   string                    `json:"phone_number" valid:"required"`
	Status        string                    `json:"status"`
	CreatedAt     time.Time                 `json:"created_at"`
	CreatedBy     readUserWarehousePayload  `json:"created_by"`
	UpdatedAt     *time.Time                `json:"updated_at"`
	UpdatedBy     *readUserWarehousePayload `json:"updated_by"`
}

func (payload *RegisterWarehousePayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *UpdateWarehousePayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *ListWarehousePayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *RegisterWarehousePayload) ToEntity(userData sqlc.GetUserBackofficeRow) (data sqlc.InsertWarehouseParams) {
	data = sqlc.InsertWarehouseParams{
		Guid:          utility.GenerateGoogleUUID(),
		WarehouseCode: payload.WarehouseCode,
		Name: sql.NullString{
			String: payload.Name,
			Valid:  true,
		},
		Address:     payload.Address,
		PhoneNumber: payload.PhoneNumber,
		CreatedBy:   userData.Guid,
	}

	return
}

func (payload *UpdateWarehousePayload) ToEntity(userData sqlc.GetUserBackofficeRow, guid string) (data sqlc.UpdateWarehouseParams) {
	data = sqlc.UpdateWarehouseParams{
		Guid: guid,
		Name: sql.NullString{
			String: payload.Name,
			Valid:  true,
		},
		Address:     payload.Address,
		PhoneNumber: payload.PhoneNumber,
		UpdatedBy: sql.NullString{
			String: userData.Guid,
			Valid:  true},
	}

	return
}

func (payload *ListWarehousePayload) ToEntity() (data sqlc.ListWarehouseParams) {
	orderParam := constants.DefaultOrderValue

	data = sqlc.ListWarehouseParams{
		SetName:          payload.Filter.SetName,
		Name:             "%" + payload.Filter.Name + "%",
		SetWarehouseCode: payload.Filter.SetWarehouseCode,
		WarehouseCode:    payload.Filter.WarehouseCode,
		SetActive:        payload.Filter.SetActive,
		Active:           payload.Filter.Active,
		LimitData:        payload.Limit,
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

func ToPayloadRegisterWarehouse(warehouseData sqlc.Warehouse, userBackoffice sqlc.GetUserBackofficeRow) (payload readRegisterWarehousePayload) {
	payload = readRegisterWarehousePayload{
		GUID:          warehouseData.Guid,
		WarehouseCode: warehouseData.WarehouseCode,
		Name:          warehouseData.Name.String,
		Address:       warehouseData.Address,
		PhoneNumber:   warehouseData.PhoneNumber,
		CreatedAt:     warehouseData.CreatedAt,
		CreatedBy: readUserWarehousePayload{
			GUID: userBackoffice.Guid,
		},
	}

	if userBackoffice.Name.Valid {
		payload.CreatedBy.Name = userBackoffice.Name.String
	}

	if warehouseData.DeletedAt.Valid {
		payload.Status = constants.StatusInactive
	} else {
		payload.Status = constants.StatusActive
	}

	return
}

func ToPayloadUpdateWarehouse(warehouseData sqlc.Warehouse, userBackoffice sqlc.GetUserBackofficeRow) (payload readUpdateWarehousePayload) {
	payload = readUpdateWarehousePayload{
		GUID:          warehouseData.Guid,
		WarehouseCode: warehouseData.WarehouseCode,
		Name:          warehouseData.Name.String,
		Address:       warehouseData.Address,
		PhoneNumber:   warehouseData.PhoneNumber,
		UpdatedAt:     warehouseData.UpdatedAt.Time,
		UpdatedBy: readUserWarehousePayload{
			GUID: userBackoffice.Guid,
		},
	}

	if userBackoffice.Name.Valid {
		payload.UpdatedBy.Name = userBackoffice.Name.String
	}

	if warehouseData.DeletedAt.Valid {
		payload.Status = constants.StatusInactive
	} else {
		payload.Status = constants.StatusActive
	}

	return
}

func ToPayloadWarehouse(warehouseData sqlc.GetWarehouseRow) (payload readWarehousePayload) {
	payload = readWarehousePayload{
		GUID:          warehouseData.Guid,
		WarehouseCode: warehouseData.WarehouseCode,
		Name:          warehouseData.Name.String,
		Address:       warehouseData.Address,
		PhoneNumber:   warehouseData.PhoneNumber,
		CreatedAt:     warehouseData.CreatedAt,
		CreatedBy: readUserWarehousePayload{
			GUID: warehouseData.UserID.String,
		},
	}

	if warehouseData.UserID.Valid {
		payload.CreatedBy.GUID = warehouseData.UserID.String
	}

	if warehouseData.UserName.Valid {
		payload.CreatedBy.Name = warehouseData.UserName.String
	}

	if warehouseData.UpdatedAt.Valid {
		payload.UpdatedAt = &warehouseData.UpdatedAt.Time
	}

	if warehouseData.UpdatedBy.Valid {
		payload.UpdatedBy = &readUserWarehousePayload{
			GUID: warehouseData.UserIDUpdate.String,
			Name: warehouseData.UserNameUpdate.String,
		}
	}

	if warehouseData.DeletedAt.Valid {
		payload.Status = constants.StatusInactive
	} else {
		payload.Status = constants.StatusActive
	}

	return
}

func ToPayloadListWarehouse(listWarehouse []sqlc.ListWarehouseRow) (payload []*readWarehousePayload) {
	payload = make([]*readWarehousePayload, len(listWarehouse))

	for i := range listWarehouse {
		payload[i] = new(readWarehousePayload)
		data := ToPayloadWarehouse(sqlc.GetWarehouseRow(listWarehouse[i]))
		payload[i] = &data
	}

	return
}
