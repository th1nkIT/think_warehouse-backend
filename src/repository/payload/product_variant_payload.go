package payload

import "time"

type readProductVariant struct {
	ID              int64                 `json:"-"`
	GUID            string                `json:"id"`
	ProductID       string                `json:"product_id"`
	Name            string                `json:"name"`
	Sku             string                `json:"sku"`
	IsActive        bool                  `json:"is_active"`
	CreatedAt       time.Time             `json:"created_at"`
	CreatedBy       string                `json:"created_by"`
	UpdatedAt       *time.Time            `json:"updated_at"`
	UpdatedBy       *string               `json:"updated_by"`
	ProductPrice    readProductPrice      `json:"product_price"`
	ProductStock    readProductStock      `json:"product_stock"`
	ProductStockLog []readProductStockLog `json:"product_stock_log"`
}
