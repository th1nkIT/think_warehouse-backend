package payload

import (
	"database/sql"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"think_warehouse/common/constants"
	"think_warehouse/common/httpservice"
	"think_warehouse/common/utility"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"time"
)

type RegisterProductCategoryPayload struct {
	Name string `json:"name" valid:"required"`
}

type UpdateProductCategoryPayload struct {
	Name string `json:"name" valid:"required"`
}

type ListProductCategoryPayload struct {
	Filter ListProductCategoryFilterPayload `json:"filter"`
	Limit  int32                            `json:"limit" valid:"required"`
	Offset int32                            `json:"page" valid:"required"`
	Order  string                           `json:"order" valid:"required"`
	Sort   string                           `json:"sort" valid:"required"` // ASC, DESC
}

type ListProductCategoryFilterPayload struct {
	SetName   bool   `json:"set_name"`
	Name      string `json:"name"`
	SetActive bool   `json:"set_active"`
	Active    string `json:"active"`
}

type readRegisterProductCategoryPayload struct {
	GUID      string                    `json:"id"`
	Name      string                    `json:"name"`
	Status    string                    `json:"status"`
	CreatedAt time.Time                 `json:"created_at"`
	CreatedBy readUserBackofficePayload `json:"created_by"`
}

type readUpdateProductCategoryPayload struct {
	GUID      string                    `json:"id"`
	Name      string                    `json:"name"`
	Status    string                    `json:"status"`
	UpdatedAt time.Time                 `json:"updated_at"`
	UpdatedBy readUserBackofficePayload `json:"updated_by"`
}

type readProductCategoryPayload struct {
	GUID      string                     `json:"id"`
	Name      string                     `json:"name"`
	Status    string                     `json:"status"`
	CreatedAt time.Time                  `json:"created_at"`
	CreatedBy readUserBackofficePayload  `json:"created_by"`
	UpdatedAt *time.Time                 `json:"updated_at"`
	UpdatedBy *readUserBackofficePayload `json:"updated_by"`
	DeletedAt *time.Time                 `json:"deleted_at"`
	DeletedBy *readUserBackofficePayload `json:"deleted_by"`
}

type readProductCategory struct {
	ID        int64      `json:"-"`
	GUID      string     `json:"id"`
	Name      string     `json:"name"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *string    `json:"updated_by"`
}

func (payload *RegisterProductCategoryPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *UpdateProductCategoryPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *ListProductCategoryPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *RegisterProductCategoryPayload) ToEntity(userData sqlc.GetUserBackofficeRow) (data sqlc.InsertProductCategoryParams) {
	data = sqlc.InsertProductCategoryParams{
		Guid:      utility.GenerateGoogleUUID(),
		Name:      payload.Name,
		CreatedBy: userData.Guid,
	}

	return
}

func (payload *UpdateProductCategoryPayload) ToEntity(userData sqlc.GetUserBackofficeRow, guid string) (data sqlc.UpdateProductCategoryParams) {
	data = sqlc.UpdateProductCategoryParams{
		Name: payload.Name,
		UpdatedBy: sql.NullString{
			String: userData.Guid,
			Valid:  true,
		},
		Guid: guid,
	}

	return
}

func (payload *ListProductCategoryPayload) ToEntity() (data sqlc.ListProductCategoryParams) {
	orderParam := constants.DefaultOrderValue

	data = sqlc.ListProductCategoryParams{
		SetName:   payload.Filter.SetName,
		Name:      "%" + payload.Filter.Name + "%",
		SetActive: payload.Filter.SetActive,
		Active:    payload.Filter.Active,
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

func ToPayloadRegisterProductCategory(productCategoryData sqlc.ProductCategory, userData sqlc.GetUserBackofficeRow) (payload readRegisterProductCategoryPayload) {
	payload = readRegisterProductCategoryPayload{
		GUID:      productCategoryData.Guid,
		Name:      productCategoryData.Name,
		CreatedAt: productCategoryData.CreatedAt,
		CreatedBy: readUserBackofficePayload{
			GUID: userData.Guid,
		},
	}

	if userData.Name.Valid {
		payload.CreatedBy.Name = userData.Name.String
	}

	if !productCategoryData.DeletedAt.Valid {
		payload.Status = constants.StatusActive
	} else {
		payload.Status = constants.StatusInactive
	}

	return
}

func ToPayloadUpdateProductCategory(productCategoryData sqlc.ProductCategory, userData sqlc.GetUserBackofficeRow) (payload readUpdateProductCategoryPayload) {
	payload = readUpdateProductCategoryPayload{
		GUID:      productCategoryData.Guid,
		Name:      productCategoryData.Name,
		UpdatedAt: productCategoryData.UpdatedAt.Time,
		UpdatedBy: readUserBackofficePayload{
			GUID: userData.Guid,
		},
	}

	if userData.Name.Valid {
		payload.UpdatedBy.Name = userData.Name.String
	}

	if productCategoryData.DeletedAt.Valid {
		payload.Status = constants.StatusInactive
	} else {
		payload.Status = constants.StatusActive
	}

	return
}

func ToPayloadProductCategory(productCategoryData sqlc.GetProductCategoryRow) (payload readProductCategoryPayload) {
	payload = readProductCategoryPayload{
		GUID:      productCategoryData.Guid,
		Name:      productCategoryData.Name,
		CreatedAt: productCategoryData.CreatedAt,
	}

	if productCategoryData.UserID.Valid {
		payload.CreatedBy.GUID = productCategoryData.UserID.String
	}

	if productCategoryData.UserName.Valid {
		payload.CreatedBy.Name = productCategoryData.UserName.String
	}

	if productCategoryData.UpdatedAt.Valid {
		payload.UpdatedAt = &productCategoryData.UpdatedAt.Time
	}

	if productCategoryData.UpdatedBy.Valid {
		payload.UpdatedBy = &readUserBackofficePayload{
			GUID: productCategoryData.UserIDUpdate.String,
			Name: productCategoryData.UserNameUpdate.String,
		}
	}

	if productCategoryData.DeletedAt.Valid {
		payload.Status = constants.StatusInactive
	} else {
		payload.Status = constants.StatusActive
	}

	if productCategoryData.DeletedAt.Valid {
		payload.DeletedAt = &productCategoryData.DeletedAt.Time
	}

	if productCategoryData.DeletedBy.Valid {
		payload.DeletedBy = &readUserBackofficePayload{
			GUID: productCategoryData.UserIDDelete.String,
			Name: productCategoryData.UserNameDelete.String,
		}
	}

	return
}

func ToPayloadListProductCategory(listProductCategory []sqlc.ListProductCategoryRow) (payload []*readProductCategoryPayload) {
	payload = make([]*readProductCategoryPayload, len(listProductCategory))

	for i := range listProductCategory {
		payload[i] = new(readProductCategoryPayload)
		data := ToPayloadProductCategory(sqlc.GetProductCategoryRow(listProductCategory[i]))
		payload[i] = &data
	}

	return
}
