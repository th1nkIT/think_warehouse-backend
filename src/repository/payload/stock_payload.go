package payload

import (
	"database/sql"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"think_warehouse/common/constants"
	"think_warehouse/common/httpservice"
	"think_warehouse/common/utility"
	sqlc "think_warehouse/src/repository/pgbo_sqlc"
)

type CreateStockPayload struct {
	Stock            int64  `json:"stock"`
	ProductID        string `json:"product_id" valid:"required"`
	ProductVariantID string `json:"product_variant_id"`
}

type readProductStocks struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	SKU   string `json:"sku"`
}

type readStock struct {
	GUID           string            `json:"id"`
	Stock          int64             `json:"stock"`
	StockRemaining int64             `json:"stock_remaining"`
	Product        readProductStocks `json:"product"`
	ProductVariant readProductStocks `json:"product_variant"`
}

type ListStockPayload struct {
	Filter ListStockFilterPayload `json:"filter"`
	Limit  int32                  `json:"limit" valid:"required"`
	Offset int32                  `json:"page" valid:"required"`
	Order  string                 `json:"order" valid:"required"`
	Sort   string                 `json:"sort" valid:"required"` // ASC, DESC
}

type ListStockFilterPayload struct {
	SetStockGreater bool   `json:"set_stock_greater"`
	Stock           int64  `json:"stock"`
	SetStockLower   bool   `json:"set_stock_lower"`
	SetProductName  bool   `json:"set_product_name"`
	ProductName     string `json:"product_name"`
	SetVariantName  bool   `json:"set_variant_name"`
	VariantName     string `json:"variant_name"`
}

type UpdateStockPayload struct {
	GUID      string `json:"id" valid:"required"`
	Stock     int64  `json:"stock"`
	StockType sqlc.StockTypeEnum
}

type GetStockByProductAndVariantParams struct {
	ProductID           string `json:"product_id"`
	ProductVariantID    string `json:"product_variant_id"`
	SetProductVariantID bool   `json:"set_product_variant_id"`
}

func (payload *CreateStockPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *UpdateStockPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *ListStockPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *UpdateStockPayload) ToEntityUpdateStock(userData sqlc.GetUserBackofficeRow, stock sqlc.GetStockRow) (updateStock sqlc.UpdateStockParams, insertLog sqlc.InsertStockLogParams) {
	// STOCK
	updateStock = sqlc.UpdateStockParams{
		Guid: stock.Guid,
		UpdatedBy: sql.NullString{
			String: userData.Guid,
			Valid:  true,
		},
		Stock: payload.Stock,
	}

	// STOCK LOG
	insertLog.Guid = utility.GenerateGoogleUUID()
	insertLog.Note = sql.NullString{String: "Update stock", Valid: true}
	insertLog.ProductID = sql.NullString{String: stock.ProductID.String, Valid: true}

	if stock.ProductVariantID.Valid {
		insertLog.ProductVariantID = sql.NullString{String: stock.ProductVariantID.String, Valid: true}
	}

	if payload.Stock > stock.Stock {
		insertLog.StockType = sqlc.NullStockTypeEnum{
			Valid:         true,
			StockTypeEnum: sqlc.StockTypeEnumOUT,
		}
		insertLog.StockLog = int32(payload.Stock - stock.Stock)
	} else {
		insertLog.StockType = sqlc.NullStockTypeEnum{
			Valid:         true,
			StockTypeEnum: sqlc.StockTypeEnumIN,
		}
		insertLog.StockLog = int32(stock.Stock - payload.Stock)
	}

	return
}

func (payload *GetStockByProductAndVariantParams) ToEntity() (data sqlc.GetStockByProductAndVariantParams) {
	data = sqlc.GetStockByProductAndVariantParams{
		SetProductVariantID: payload.SetProductVariantID,
	}

	if payload.ProductID != "" {
		data.ProductID.Valid = true
		data.ProductID.String = payload.ProductID
	}

	if payload.ProductVariantID != "" {
		data.ProductVariantID.Valid = true
		data.ProductVariantID.String = payload.ProductVariantID
	}

	return
}

func (payload *ListStockPayload) ToEntity() (data sqlc.ListStockParams) {
	orderParam := constants.DefaultOrderValue

	data = sqlc.ListStockParams{
		SetStockGreater: payload.Filter.SetStockGreater,
		SetStockLower:   payload.Filter.SetStockLower,
		Stock:           payload.Filter.Stock,
		SetProductName:  payload.Filter.SetProductName,
		ProductName:     "%" + payload.Filter.ProductName + "%",
		SetVariantName:  payload.Filter.SetVariantName,
		VariantName:     "%" + payload.Filter.VariantName + "%",
		LimitData:       payload.Limit,
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

func ToPayloadUpdateStock(stock []sqlc.Stock) (response []readStock) {
	for i := range stock {
		out := readStock{
			GUID:           stock[i].Guid,
			Stock:          stock[i].Stock,
			StockRemaining: stock[i].Stock,
			Product: readProductStocks{
				ID:    "",
				Title: "",
			},
			ProductVariant: readProductStocks{
				ID:    "",
				Title: "",
			},
		}

		if stock[i].ProductID.Valid {
			out.Product.ID = stock[i].ProductID.String
		}

		if stock[i].ProductVariantID.Valid {
			out.ProductVariant.ID = stock[i].ProductVariantID.String
		}

		response = append(response, out)
	}

	return
}

func ToPayloadStockSingle(stock sqlc.GetStockRow) (payload readStock) {
	payload = readStock{
		GUID:           stock.Guid,
		Stock:          stock.Stock,
		StockRemaining: stock.Stock,
		Product: readProductStocks{
			ID:    stock.ProductID.String,
			Title: stock.ProductName.String,
			SKU:   stock.ProductSku.String,
		},
		ProductVariant: readProductStocks{
			ID:    stock.ProductVariantID.String,
			Title: stock.VariantName.String,
			SKU:   stock.VariantSku.String,
		},
	}

	if stock.ProductID.Valid {
		payload.Product.ID = stock.ProductID.String
		if stock.ProductName.Valid {
			payload.Product.Title = stock.ProductName.String
		}
		if stock.ProductSku.Valid {
			payload.Product.SKU = stock.ProductSku.String
		}
	}

	if stock.ProductVariantID.Valid {
		payload.ProductVariant.ID = stock.ProductVariantID.String
	}

	return
}

func ToPayloadStockArray(stock sqlc.ListStockRow) (payload readStock) {
	payload = readStock{
		GUID:           stock.Guid,
		Stock:          stock.Stock,
		StockRemaining: stock.Stock,
		Product: readProductStocks{
			ID:    stock.ProductID.String,
			Title: stock.ProductName.String,
			SKU:   stock.ProductSku,
		},
		ProductVariant: readProductStocks{
			ID:    stock.ProductVariantID.String,
			Title: stock.VariantName.String,
			SKU:   stock.VariantSku.String,
		},
	}

	if stock.ProductVariantID.Valid {
		payload.ProductVariant.ID = stock.ProductVariantID.String
		if stock.VariantName.Valid {
			payload.ProductVariant.Title = stock.VariantName.String
		}

		if stock.VariantSku.Valid {
			payload.ProductVariant.SKU = stock.VariantSku.String
		}
	}

	return
}

func ToPayloadListStock(listStock []sqlc.ListStockRow) (payload []*readStock) {
	payload = make([]*readStock, len(listStock))

	for i := range listStock {
		payload[i] = new(readStock)
		data := ToPayloadStockArray(listStock[i])
		payload[i] = &data
	}

	return
}
