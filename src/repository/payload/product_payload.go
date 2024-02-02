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

type CreateProductWithoutVariantRequest struct {
	ProductPictureUrl string `json:"profile_picture_url"`
	IsVariant         bool   `json:"is_variant"`
	ProductCategoryID string `json:"product_category_id" valid:"required~product_category_id is required field"`
	Name              string `json:"name" valid:"required~name is required field"`
	ProductCode       string `json:"product_code"`
	Description       string `json:"description"`
	ProductSKU        string `json:"product_sku"`
	Price             int64  `json:"price"`
	DiscountType      string `json:"discount_type"`
	Discount          int64  `json:"discount"`
	Stock             int64  `json:"stock"`
	IsActive          bool   `json:"is_active" valid:"required~is_active is required field"`
}

type CreateProductWithVariantRequest struct {
	ProductPictureUrl string                 `json:"profile_picture_url"`
	ProductCode       string                 `json:"product_code"`
	IsVariant         bool                   `json:"is_variant"`
	ProductCategoryID string                 `json:"product_category_id" valid:"required~product_category_id is required field"`
	Name              string                 `json:"name" valid:"required~name is required field"`
	ProductSKU        string                 `json:"product_sku"`
	Description       string                 `json:"description"`
	ProductVariant    []CreateProductVariant `json:"product_variant"`
}

type CreateProductVariant struct {
	GUID         string `json:"id"`
	Name         string `json:"name" valid:"required~variant_name is required field"`
	Sku          string `json:"sku"`
	Price        int64  `json:"price"`
	DiscountType string `json:"discount_type"`
	Discount     int64  `json:"discount"`
	Stock        int64  `json:"stock"`
	IsActive     bool   `json:"is_active" valid:"required~is_active is required field"`
}

type ProductWithoutVariant struct {
	ID                int64                 `json:"-"`
	GUID              string                `json:"id"`
	ImageURL          *string               `json:"image_url"`
	IsVariant         bool                  `json:"is_variant"`
	ProductCategoryID string                `json:"product_category_id,omitempty"`
	ProductCode       string                `json:"product_code"`
	Name              string                `json:"name"`
	Description       *string               `json:"description"`
	Sku               *string               `json:"sku"`
	IsActive          bool                  `json:"is_active"`
	CreatedAt         time.Time             `json:"created_at"`
	CreatedBy         string                `json:"created_by"`
	UpdatedAt         *time.Time            `json:"updated_at"`
	UpdatedBy         *string               `json:"updated_by"`
	Duration          *int32                `json:"duration"`
	ProductCategory   readProductCategory   `json:"product_category,omitempty"`
	ProductPrice      readProductPrice      `json:"product_price"`
	ProductStockLog   []readProductStockLog `json:"product_stock_log"`
	ProductStock      readProductStock      `json:"product_stock"`
}

type ProductWithVariant struct {
	ID                int64                `json:"-"`
	GUID              string               `json:"id"`
	ImageURL          *string              `json:"image_url"`
	IsVariant         bool                 `json:"is_variant"`
	ProductCategoryID string               `json:"product_category_id,omitempty"`
	Name              string               `json:"name"`
	ProductCode       string               `json:"product_code"`
	IsActive          bool                 `json:"is_active"`
	CreatedAt         time.Time            `json:"created_at"`
	CreatedBy         string               `json:"created_by"`
	UpdatedAt         *time.Time           `json:"updated_at"`
	UpdatedBy         *string              `json:"updated_by"`
	ProductCategory   readProductCategory  `json:"product_category,omitempty"`
	ProductVariant    []readProductVariant `json:"product_variant"`
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

type readUpdateProductPayload struct {
	GUID              string                    `json:"id"`
	Name              string                    `json:"name"`
	ProductCode       string                    `json:"product_code"`
	ProductPictureUrl *string                   `json:"profile_picture_image_url"`
	Description       string                    `json:"description"`
	Status            string                    `json:"status"`
	CategoryId        string                    `json:"category_id"`
	Category          readProductCategory       `json:"category"`
	ProductSKU        string                    `json:"product_sku"`
	IsVariant         bool                      `json:"is_variant"`
	UpdatedAt         time.Time                 `json:"updated_at"`
	UpdatedBy         readUserBackofficePayload `json:"updated_by"`
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
	Category          readProductCategory        `json:"category"`
	CreatedAt         time.Time                  `json:"created_at"`
	CreatedBy         readUserBackofficePayload  `json:"created_by"`
	UpdatedAt         *time.Time                 `json:"updated_at"`
	UpdatedBy         *readUserBackofficePayload `json:"updated_by"`
	ProductVariant    json.RawMessage            `json:"product_variant"`
}

func (payload *CreateProductWithoutVariantRequest) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *CreateProductWithVariantRequest) Validate() (err error) {
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

func (payload *CreateProductWithVariantRequest) ToEntity(userData sqlc.GetUserBackofficeRow) (
	productParams sqlc.InsertProductParams,
	productVariantParams []sqlc.InsertProductVariantParams,
	productPriceParams []sqlc.InsertProductPriceParams,
	stockLogParams []sqlc.InsertStockLogParams,
	stockParams []sqlc.InsertStockParams,
) {
	productParams = sqlc.InsertProductParams{
		Guid:              utility.GenerateGoogleUUID(),
		ProductPictureUrl: sql.NullString{String: payload.ProductPictureUrl, Valid: len(payload.ProductPictureUrl) > 0},
		IsVariant:         true,
		CategoryID:        payload.ProductCategoryID,
		Name: sql.NullString{
			String: payload.Name,
			Valid:  true,
		},
		ProductSku:  payload.ProductSKU,
		Description: payload.Description,
		ProductCode: payload.ProductCode,
		CreatedBy:   userData.Guid,
	}

	for i := range payload.ProductVariant {
		productVariant := payload.ProductVariant[i]
		productVariantParam := sqlc.InsertProductVariantParams{
			Guid:      utility.GenerateGoogleUUID(),
			ProductID: productParams.Guid,
			Name:      productVariant.Name,
			Sku:       productVariant.Sku,
			IsActive:  productVariant.IsActive,
			CreatedBy: userData.Guid,
		}

		productPriceParam := sqlc.InsertProductPriceParams{
			Guid:             utility.GenerateGoogleUUID(),
			ProductID:        sql.NullString{String: productParams.Guid, Valid: true},
			ProductVariantID: sql.NullString{String: productVariantParam.Guid, Valid: true},
			Price:            productVariant.Price,
			DiscountType:     sqlc.DiscountTypeEnum(productVariant.DiscountType),
			Discount:         sql.NullInt64{Int64: productVariant.Discount, Valid: productVariant.Discount > 0},
			IsActive:         true,
			CreatedBy:        userData.Guid,
		}

		stockLogParam := sqlc.InsertStockLogParams{
			Guid:             utility.GenerateGoogleUUID(),
			ProductID:        sql.NullString{String: productParams.Guid, Valid: true},
			ProductVariantID: sql.NullString{String: productVariantParam.Guid, Valid: true},
			StockLog:         int32(productVariant.Stock),
			StockType: sqlc.NullStockTypeEnum{
				StockTypeEnum: sqlc.StockTypeEnumIN,
			},
			Note:      sql.NullString{String: "Create stock", Valid: true},
			CreatedBy: userData.Guid,
		}

		stockParam := sqlc.InsertStockParams{
			Guid:             utility.GenerateGoogleUUID(),
			ProductID:        sql.NullString{String: productParams.Guid, Valid: true},
			ProductVariantID: sql.NullString{String: productVariantParam.Guid, Valid: true},
			Stock:            productVariant.Stock,
			CreatedBy:        userData.Guid,
		}

		productVariantParams = append(productVariantParams, productVariantParam)
		productPriceParams = append(productPriceParams, productPriceParam)
		stockLogParams = append(stockLogParams, stockLogParam)
		stockParams = append(stockParams, stockParam)
	}

	return
}

func (payload *CreateProductWithoutVariantRequest) ToEntity(userData sqlc.GetUserBackofficeRow) (
	product sqlc.InsertProductParams,
	productPrice sqlc.InsertProductPriceParams,
	stockLog sqlc.InsertStockLogParams,
	stock sqlc.InsertStockParams,
) {
	product = sqlc.InsertProductParams{
		Guid: utility.GenerateGoogleUUID(),
		Name: sql.NullString{
			String: payload.Name,
			Valid:  true,
		},
		CategoryID:  payload.ProductCategoryID,
		ProductSku:  payload.ProductSKU,
		IsVariant:   payload.IsVariant,
		ProductCode: payload.ProductCode,
		ProductPictureUrl: sql.NullString{
			String: payload.ProductPictureUrl,
			Valid:  true,
		},
		Description: payload.Description,
		CreatedBy:   userData.Guid,
	}

	productPrice = sqlc.InsertProductPriceParams{
		Guid: utility.GenerateGoogleUUID(),
		ProductID: sql.NullString{
			String: product.Guid,
			Valid:  true,
		},
		Price:        payload.Price,
		DiscountType: sqlc.DiscountTypeEnum(payload.DiscountType),
		Discount: sql.NullInt64{
			Int64: payload.Discount,
			Valid: true,
		},
		IsActive:  true,
		CreatedBy: userData.Guid,
	}

	stockLog = sqlc.InsertStockLogParams{
		Guid: utility.GenerateGoogleUUID(),
		ProductID: sql.NullString{
			String: product.Guid,
			Valid:  true,
		},
		StockLog: int32(payload.Stock),
		StockType: sqlc.NullStockTypeEnum{
			StockTypeEnum: sqlc.StockTypeEnumIN,
			Valid:         true,
		},
		Note:      sql.NullString{String: "Create stock", Valid: true},
		CreatedBy: userData.Guid,
	}

	stock = sqlc.InsertStockParams{
		Guid: utility.GenerateGoogleUUID(),
		ProductID: sql.NullString{
			String: product.Guid,
			Valid:  true,
		},
		Stock:     payload.Stock,
		CreatedBy: userData.Guid,
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

func ToResponsePayloadProductWithVariant(product sqlc.Product, productCategory sqlc.GetProductCategoryRow, productVariants []sqlc.ProductVariant, productPrices []sqlc.ProductPrice, stockLogs []sqlc.StockLog, stocks []sqlc.Stock) (payload ProductWithVariant) {
	payload = ProductWithVariant{
		ID:                product.ID,
		GUID:              product.Guid,
		IsVariant:         product.IsVariant,
		ImageURL:          &product.ProductPictureUrl.String,
		ProductCategoryID: product.CategoryID,
		ProductCode:       product.ProductCode,
		Name:              product.Name.String,
		CreatedAt:         product.CreatedAt,
		CreatedBy:         product.CreatedBy,
	}

	if product.DeletedAt.Valid {
		payload.IsActive = false
	} else {
		payload.IsActive = true
	}

	if product.UpdatedAt.Valid {
		payload.UpdatedAt = &product.UpdatedAt.Time
	}

	if product.UpdatedBy.Valid {
		payload.UpdatedBy = &product.UpdatedBy.String
	}

	payload.ProductCategory = readProductCategory{
		GUID:      productCategory.Guid,
		Name:      productCategory.Name,
		CreatedAt: productCategory.CreatedAt,
		CreatedBy: productCategory.CreatedBy,
	}

	if productCategory.UpdatedAt.Valid {
		payload.ProductCategory.UpdatedAt = &productCategory.UpdatedAt.Time
	}

	if productCategory.UpdatedBy.Valid {
		payload.ProductCategory.UpdatedBy = &productCategory.UpdatedBy.String
	}

	for i := range productVariants {
		var (
			productVariantPayload readProductVariant
			productVariant        sqlc.ProductVariant
			productPrice          sqlc.ProductPrice
			stockLog              sqlc.StockLog
			stock                 sqlc.Stock
		)

		productVariant = productVariants[i]
		productPrice = productPrices[i]
		stockLog = stockLogs[i]
		stock = stocks[i]

		productVariantPayload = readProductVariant{
			ID:        productVariant.ID,
			GUID:      productVariant.Guid,
			ProductID: productVariant.ProductID,
			Name:      productVariant.Name,
			Sku:       productVariant.Sku,
			IsActive:  productVariant.IsActive,
			CreatedAt: productVariant.CreatedAt,
			CreatedBy: productVariant.CreatedBy,
		}

		if productVariant.UpdatedAt.Valid {
			productVariantPayload.UpdatedAt = &productVariant.UpdatedAt.Time
		}

		if productVariant.UpdatedBy.Valid {
			productVariantPayload.UpdatedBy = &productVariant.UpdatedBy.String
		}

		productVariantPayload.ProductPrice = readProductPrice{
			ID:        productPrice.ID,
			GUID:      productPrice.Guid,
			Price:     productPrice.Price,
			IsActive:  productPrice.IsActive,
			CreatedAt: productPrice.CreatedAt,
			CreatedBy: productPrice.CreatedBy,
		}

		if productPrice.DiscountType.Valid {
			productVariantPayload.ProductPrice.DiscountType = (*string)(&productPrice.DiscountType.DiscountTypeEnum)
		}

		if productPrice.Discount.Valid {
			productVariantPayload.ProductPrice.Discount = &productPrice.Discount.Int64
		}

		if productPrice.UpdatedAt.Valid {
			productVariantPayload.ProductPrice.UpdatedAt = &productPrice.UpdatedAt.Time
		}

		if productPrice.UpdatedBy.Valid {
			productVariantPayload.ProductPrice.UpdatedBy = &productPrice.UpdatedBy.String
		}

		productVariantPayload.ProductStock = readProductStock{
			ID:        stock.ID,
			GUID:      stock.Guid,
			Stock:     stock.Stock,
			CreatedAt: stock.CreatedAt,
			CreatedBy: stock.CreatedBy,
		}

		stockLogSing := readProductStockLog{
			ID:        stockLog.ID,
			GUID:      stockLog.Guid,
			StockLog:  int64(stockLog.StockLog),
			StockType: string(stockLog.StockType.StockTypeEnum),
			CreatedAt: stockLog.CreatedAt,
			CreatedBy: stockLog.CreatedBy,
		}

		productVariantPayload.ProductStockLog = append(productVariantPayload.ProductStockLog, stockLogSing)

		if stockLog.UpdatedAt.Valid {
			productVariantPayload.ProductStock.UpdatedAt = &stockLog.UpdatedAt.Time
		}

		if stockLog.UpdatedBy.Valid {
			productVariantPayload.ProductStock.UpdatedBy = &stockLog.UpdatedBy.String
		}

		payload.ProductVariant = append(payload.ProductVariant, productVariantPayload)
	}

	return
}

func ToResponsePayloadProductWithoutVariant(product sqlc.Product, productCategory sqlc.GetProductCategoryRow, productPrice sqlc.ProductPrice, stockLog sqlc.StockLog, stock sqlc.Stock) (payload ProductWithoutVariant) {
	payload = ProductWithoutVariant{
		ID:                product.ID,
		GUID:              product.Guid,
		IsVariant:         product.IsVariant,
		ProductCategoryID: product.CategoryID,
		ProductCode:       product.ProductCode,
		Name:              product.Name.String,
		Description:       &product.Description,
		Sku:               &product.ProductSku,
		CreatedAt:         product.CreatedAt,
		CreatedBy:         product.CreatedBy,
	}

	if product.DeletedAt.Valid {
		payload.IsActive = false
	} else {
		payload.IsActive = true
	}

	if product.ProductPictureUrl.Valid {
		payload.ImageURL = &product.ProductPictureUrl.String
	}

	if product.UpdatedAt.Valid {
		payload.UpdatedAt = &product.UpdatedAt.Time
	}

	if product.UpdatedBy.Valid {
		payload.UpdatedBy = &product.UpdatedBy.String
	}

	payload.ProductCategory = readProductCategory{
		GUID:      productCategory.Guid,
		Name:      productCategory.Name,
		CreatedAt: productCategory.CreatedAt,
		CreatedBy: productCategory.CreatedBy,
	}

	if productCategory.UpdatedAt.Valid {
		payload.ProductCategory.UpdatedAt = &productCategory.UpdatedAt.Time
	}

	if productCategory.UpdatedBy.Valid {
		payload.ProductCategory.UpdatedBy = &productCategory.UpdatedBy.String
	}

	payload.ProductPrice = readProductPrice{
		ID:        productPrice.ID,
		GUID:      productPrice.Guid,
		Price:     productPrice.Price,
		IsActive:  productPrice.IsActive,
		CreatedAt: productPrice.CreatedAt,
		CreatedBy: productPrice.CreatedBy,
	}

	if productPrice.DiscountType.Valid {
		payload.ProductPrice.DiscountType = (*string)(&productPrice.DiscountType.DiscountTypeEnum)
	}

	if productPrice.Discount.Valid {
		payload.ProductPrice.Discount = &productPrice.Discount.Int64
	}

	if productPrice.UpdatedAt.Valid {
		payload.ProductPrice.UpdatedAt = &productPrice.UpdatedAt.Time
	}

	if productPrice.UpdatedBy.Valid {
		payload.ProductPrice.UpdatedBy = &productPrice.UpdatedBy.String
	}

	stockLogSing := readProductStockLog{
		ID:        stockLog.ID,
		GUID:      stockLog.Guid,
		StockLog:  int64(stockLog.StockLog),
		StockType: string(stockLog.StockType.StockTypeEnum),
		CreatedAt: stockLog.CreatedAt,
		CreatedBy: stockLog.CreatedBy,
	}

	payload.ProductStockLog = append(payload.ProductStockLog, stockLogSing)

	payload.ProductStock = readProductStock{
		ID:        stock.ID,
		GUID:      stock.Guid,
		Stock:     stock.Stock,
		CreatedAt: stock.CreatedAt,
		CreatedBy: stock.CreatedBy,
	}

	if stockLog.UpdatedAt.Valid {
		payload.ProductStock.UpdatedAt = &stockLog.UpdatedAt.Time
	}

	if stockLog.UpdatedBy.Valid {
		payload.ProductStock.UpdatedBy = &stockLog.UpdatedBy.String
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
		Category: readProductCategory{
			GUID: categoryData.Guid,
			Name: categoryData.Name,
		},
		UpdatedAt: productData.UpdatedAt.Time,
		UpdatedBy: readUserBackofficePayload{
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
		Category: readProductCategory{
			GUID: productData.CategoryID,
			Name: productData.ProductCategoryName.String,
		},
		CreatedAt: productData.CreatedAt,
		CreatedBy: readUserBackofficePayload{
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
		payload.UpdatedBy = &readUserBackofficePayload{
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
