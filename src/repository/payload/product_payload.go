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

type RegisterProductPayload struct {
	Name              string `json:"name" valid:"required"`
	ProductPictureUrl string `json:"profile_picture_url"`
	Description       string `json:"description"`
}

type UpdateProductPayload struct {
	Name              string `json:"name" valid:"required"`
	ProductPictureUrl string `json:"profile_picture_url"`
	Description       string `json:"description"`
}

type ListProductPayload struct {
	Filter ListProductFilterPayload `json:"filter"`
	Limit  int32                    `json:"limit" valid:"required"`
	Offset int32                    `json:"page" valid:"required"`
	Order  string                   `json:"order" valid:"required"`
	Sort   string                   `json:"sort" valid:"required"` // ASC, DESC
}

type ListProductFilterPayload struct {
	SetName bool   `json:"set_name"`
	Name    string `json:"name"`
}

type readRegisterProductPayload struct {
	GUID              string                    `json:"id"`
	Name              string                    `json:"name"`
	ProductPictureUrl *string                   `json:"profile_picture_image_url"`
	Description       string                    `json:"description"`
	Status            string                    `json:"status"`
	CreatedAt         time.Time                 `json:"created_at"`
	CreatedBy         readUserBackOfficePayload `json:"created_by"`
}

type readUserBackOfficePayload struct {
	GUID string `json:"id"`
	Name string `json:"name"`
}

type readUpdateProductPayload struct {
	GUID              string                    `json:"id"`
	Name              string                    `json:"name"`
	ProductPictureUrl *string                   `json:"profile_picture_image_url"`
	Description       string                    `json:"description"`
	Status            string                    `json:"status"`
	UpdatedAt         time.Time                 `json:"updated_at"`
	UpdatedBy         readUserBackOfficePayload `json:"updated_by"`
}

type readProductPayload struct {
	GUID              string                     `json:"id"`
	Name              string                     `json:"name"`
	ProductPictureUrl *string                    `json:"profile_picture_image_url"`
	Description       string                     `json:"description"`
	Status            string                     `json:"status"`
	CreatedAt         time.Time                  `json:"created_at"`
	CreatedBy         readUserBackOfficePayload  `json:"created_by"`
	UpdatedAt         *time.Time                 `json:"updated_at"`
	UpdatedBy         *readUserBackOfficePayload `json:"updated_by"`
}

func (payload *RegisterProductPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *UpdateProductPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *ListProductPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *RegisterProductPayload) ToEntity(userData sqlc.GetUserBackofficeRow) (data sqlc.InsertProductParams) {
	data = sqlc.InsertProductParams{
		Guid: utility.GenerateGoogleUUID(),
		Name: sql.NullString{
			String: payload.Name,
			Valid:  true,
		},
		ProductPictureUrl: sql.NullString{
			String: payload.ProductPictureUrl,
			Valid:  true,
		},
		Description: payload.Description,
		CreatedBy:   userData.Guid,
	}

	return
}

func (payload *UpdateProductPayload) ToEntity(userData sqlc.GetUserBackofficeRow, guid string) (data sqlc.UpdateProductParams) {
	data = sqlc.UpdateProductParams{
		Guid: guid,
		Name: sql.NullString{
			String: payload.Name,
			Valid:  true,
		},
		ProductPictureUrl: sql.NullString{
			String: payload.ProductPictureUrl,
			Valid:  true,
		},
		Description: payload.Description,
		UpdatedBy:   sql.NullString{String: userData.Guid, Valid: true},
	}

	return
}

func (payload *ListProductPayload) ToEntity() (data sqlc.ListProductParams) {
	orderParam := constants.DefaultOrderValue

	data = sqlc.ListProductParams{
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

func ToPayloadRegisterProduct(productData sqlc.Product, userBackoffice sqlc.GetUserBackofficeRow) (payload readRegisterProductPayload) {
	payload = readRegisterProductPayload{
		GUID:        productData.Guid,
		Name:        productData.Name.String,
		Description: productData.Description,
		CreatedAt:   productData.CreatedAt,
		CreatedBy: readUserBackOfficePayload{
			GUID: userBackoffice.Guid,
		},
	}

	if userBackoffice.Name.Valid {
		payload.CreatedBy.Name = userBackoffice.Name.String
	}

	if productData.ProductPictureUrl.Valid {
		payload.ProductPictureUrl = &productData.ProductPictureUrl.String
	}

	return
}

func ToPayloadUpdateProduct(productData sqlc.Product, userBackoffice sqlc.GetUserBackofficeRow) (payload readUpdateProductPayload) {
	payload = readUpdateProductPayload{
		GUID:        productData.Guid,
		Name:        productData.Name.String,
		Description: productData.Description,
		UpdatedAt:   productData.UpdatedAt.Time,
		UpdatedBy: readUserBackOfficePayload{
			GUID: userBackoffice.Guid,
		},
	}

	if userBackoffice.Name.Valid {
		payload.UpdatedBy.Name = userBackoffice.Name.String
	}

	if productData.ProductPictureUrl.Valid {
		payload.ProductPictureUrl = &productData.ProductPictureUrl.String
	}

	if productData.DeletedAt.Valid {
		payload.Status = constants.StatusInactive
	} else {
		payload.Status = constants.StatusActive
	}

	return
}

func ToPayloadProduct(productData sqlc.GetProductRow) (payload readProductPayload) {
	payload = readProductPayload{
		GUID:        productData.Guid,
		Name:        productData.Name.String,
		Description: productData.Description,
		CreatedAt:   productData.CreatedAt,
		CreatedBy: readUserBackOfficePayload{
			GUID: productData.UserID.String,
		},
	}

	if productData.UserID.Valid {
		payload.CreatedBy.GUID = productData.UserID.String
	}

	if productData.UserName.Valid {
		payload.CreatedBy.Name = productData.UserName.String
	}

	if productData.ProductPictureUrl.Valid {
		payload.ProductPictureUrl = &productData.ProductPictureUrl.String
	}

	if productData.UpdatedAt.Valid {
		payload.UpdatedAt = &productData.UpdatedAt.Time
	}

	if productData.UpdatedBy.Valid {
		payload.UpdatedBy = &readUserBackOfficePayload{
			GUID: productData.UserIDUpdate.String,
			Name: productData.UserNameUpdate.String,
		}
	}

	if productData.DeletedAt.Valid {
		payload.Status = constants.StatusInactive
	} else {
		payload.Status = constants.StatusActive
	}

	return
}

func ToPayloadListProduct(listProduct []sqlc.ListProductRow) (payload []*readProductPayload) {
	payload = make([]*readProductPayload, len(listProduct))

	for i := range listProduct {
		payload[i] = new(readProductPayload)
		data := ToPayloadProduct(sqlc.GetProductRow(listProduct[i]))
		payload[i] = &data
	}

	return
}
