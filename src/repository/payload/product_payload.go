package payload

import (
	"database/sql"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"think_warehouse/common/constants"
	"think_warehouse/common/httpservice"
	"think_warehouse/common/utility"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
	"time"
)

type RegisterProductPayload struct {
	Name              string `json:"name" valid:"required"`
	ProductCode       string `json:"product_code" valid:"required"`
	ProductSKU        string `json:"product_sku"`
	IsVariant         bool   `json:"is_variant"`
	CategoryId        string `json:"category_id" valid:"required"`
	ProductPictureUrl string `json:"profile_picture_url"`
	Description       string `json:"description"`
}

type UpdateProductPayload struct {
	Name              string `json:"name" valid:"required"`
	CategoryId        string `json:"category_id" valid:"required"`
	ProductSKU        string `json:"product_sku"`
	IsVariant         bool   `json:"is_variant"`
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
	SetIsVariant         bool   `json:"set_is_variant"`
	IsVariant            bool   `json:"is_variant"`
	SetProductCategoryID bool   `json:"set_product_category_id"`
	ProductCategoryID    string `json:"product_category_id"`
	SetName              bool   `json:"set_name"`
	Name                 string `json:"name"`
	SetProductCode       bool   `json:"set_product_code"`
	ProductCode          string `json:"product_code"`
	SetDescription       bool   `json:"set_description"`
	Description          string `json:"description"`
	SetSku               bool   `json:"set_sku"`
	Sku                  string `json:"sku"`
	SetIsActive          bool   `json:"set_is_active"`
	IsActive             bool   `json:"is_active"`
}

type readRegisterProductPayload struct {
	GUID              string                    `json:"id"`
	Name              string                    `json:"name"`
	ProductCode       string                    `json:"product_code"`
	ProductPictureUrl *string                   `json:"profile_picture_image_url"`
	Description       string                    `json:"description"`
	CategoryId        string                    `json:"category_id"`
	Category          readProductCategoryData   `json:"category"`
	ProductSKU        string                    `json:"product_sku"`
	IsVariant         bool                      `json:"is_variant"`
	Status            string                    `json:"status"`
	CreatedAt         time.Time                 `json:"created_at"`
	CreatedBy         readUserBackOfficePayload `json:"created_by"`
}

type readUserBackOfficePayload struct {
	GUID string `json:"id"`
	Name string `json:"name"`
}

type readProductCategoryData struct {
	GUID string `json:"id"`
	Name string `json:"name"`
}

type readUpdateProductPayload struct {
	GUID              string                    `json:"id"`
	Name              string                    `json:"name"`
	ProductCode       string                    `json:"product_code"`
	ProductPictureUrl *string                   `json:"profile_picture_image_url"`
	Description       string                    `json:"description"`
	Status            string                    `json:"status"`
	CategoryId        string                    `json:"category_id"`
	Category          readProductCategoryData   `json:"category"`
	ProductSKU        string                    `json:"product_sku"`
	IsVariant         bool                      `json:"is_variant"`
	UpdatedAt         time.Time                 `json:"updated_at"`
	UpdatedBy         readUserBackOfficePayload `json:"updated_by"`
}

type readProduct struct {
	GUID              string                     `json:"id"`
	Name              string                     `json:"name"`
	ProductCode       string                     `json:"product_code"`
	ProductSKU        string                     `json:"product_sku"`
	IsVariant         bool                       `json:"is_variant"`
	ProductPictureUrl *string                    `json:"profile_picture_image_url"`
	Description       string                     `json:"description"`
	Status            string                     `json:"status"`
	CategoryId        string                     `json:"category_id"`
	Category          readProductCategoryData    `json:"category"`
	CreatedAt         time.Time                  `json:"created_at"`
	CreatedBy         readUserBackOfficePayload  `json:"created_by"`
	UpdatedAt         *time.Time                 `json:"updated_at"`
	UpdatedBy         *readUserBackOfficePayload `json:"updated_by"`
	ProductVariant    json.RawMessage            `json:"product_variant"`
}

type readProductPayload struct {
	GUID              string                     `json:"id"`
	Name              string                     `json:"name"`
	ProductCode       string                     `json:"product_code"`
	ProductSKU        string                     `json:"product_sku"`
	IsVariant         bool                       `json:"is_variant"`
	ProductPictureUrl *string                    `json:"profile_picture_image_url"`
	Description       string                     `json:"description"`
	Status            string                     `json:"status"`
	CategoryId        string                     `json:"category_id"`
	Category          readProductCategoryData    `json:"category"`
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
		ProductSku:  payload.ProductSKU,
		IsVariant:   payload.IsVariant,
		ProductCode: payload.ProductCode,
		CategoryID:  payload.CategoryId,
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
		ProductSku: payload.ProductSKU,
		IsVariant:  payload.IsVariant,
		CategoryID: payload.CategoryId,
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
		SetName:              payload.Filter.SetName,
		Name:                 "%" + payload.Filter.Name + "%",
		SetIsVariant:         payload.Filter.SetIsVariant,
		IsVariant:            payload.Filter.IsVariant,
		SetProductCategoryID: payload.Filter.SetProductCategoryID,
		ProductCategoryID:    payload.Filter.ProductCategoryID,
		SetDescription:       payload.Filter.SetDescription,
		Description:          "%" + payload.Filter.Description + "%",
		SetSku:               payload.Filter.SetSku,
		Sku:                  "%" + payload.Filter.Sku + "%",
		SetProductCode:       payload.Filter.SetProductCode,
		ProductCode:          "%" + payload.Filter.ProductCode + "%",
		LimitData:            payload.Limit,
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

func ToPayloadRegisterProduct(productData sqlc.Product, userBackoffice sqlc.GetUserBackofficeRow, categoryData sqlc.GetProductCategoryRow) (payload readRegisterProductPayload) {
	payload = readRegisterProductPayload{
		GUID:        productData.Guid,
		Name:        productData.Name.String,
		ProductCode: productData.ProductCode,
		ProductSKU:  productData.ProductSku,
		IsVariant:   productData.IsVariant,
		CategoryId:  productData.CategoryID,
		Category: readProductCategoryData{
			GUID: categoryData.Guid,
			Name: categoryData.Name,
		},
		Status:      constants.StatusActive,
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

func ToPayloadUpdateProduct(productData sqlc.Product, userBackoffice sqlc.GetUserBackofficeRow, categoryData sqlc.GetProductCategoryRow) (payload readUpdateProductPayload) {
	payload = readUpdateProductPayload{
		GUID:        productData.Guid,
		Name:        productData.Name.String,
		ProductCode: productData.ProductCode,
		ProductSKU:  productData.ProductSku,
		IsVariant:   productData.IsVariant,
		Description: productData.Description,
		CategoryId:  productData.CategoryID,
		Category: readProductCategoryData{
			GUID: categoryData.Guid,
			Name: categoryData.Name,
		},
		UpdatedAt: productData.UpdatedAt.Time,
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

// Todo: Add ToPayloadProductVariant and NoVariant

func ToPayloadProduct(productData sqlc.GetProductRow) (payload readProduct) {
	payload = readProduct{
		GUID:        productData.Guid,
		Name:        productData.Name.String,
		ProductCode: productData.ProductCode,
		IsVariant:   productData.IsVariant,
		ProductSKU:  productData.ProductSku,
		Description: productData.Description,
		CategoryId:  productData.CategoryID,
		Category: readProductCategoryData{
			GUID: productData.CategoryID,
			Name: productData.ProductCategoryName.String,
		},
		CreatedAt: productData.CreatedAt,
		CreatedBy: readUserBackOfficePayload{
			GUID: productData.CreatedByGuid.String,
		},
	}

	if productData.CreatedByGuid.Valid {
		payload.CreatedBy.GUID = productData.CreatedByGuid.String
	}

	if productData.CreatedByName.Valid {
		payload.CreatedBy.Name = productData.CreatedByName.String
	}

	if productData.ProductPictureUrl.Valid {
		payload.ProductPictureUrl = &productData.ProductPictureUrl.String
	}

	if productData.UpdatedAt.Valid {
		payload.UpdatedAt = &productData.UpdatedAt.Time
	}

	if productData.UpdatedBy.Valid {
		payload.UpdatedBy = &readUserBackOfficePayload{
			GUID: productData.UpdatedByGuid.String,
			Name: productData.UpdatedByName.String,
		}
	}

	if productData.DeletedAt.Valid {
		payload.Status = constants.StatusInactive
	} else {
		payload.Status = constants.StatusActive
	}

	if productData.ProductVariant != nil {
		payload.ProductVariant = productData.ProductVariant
	}

	return
}

func ToPayloadListProduct(listProduct []sqlc.ListProductRow) (payload []*readProduct) {
	payload = make([]*readProduct, len(listProduct))

	for i := range listProduct {
		payload[i] = new(readProduct)
		data := ToPayloadProduct(sqlc.GetProductRow(listProduct[i]))
		payload[i] = &data
	}

	return
}
