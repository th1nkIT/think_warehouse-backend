package payload

import "time"

type readProductStockLog struct {
	ID        int64      `json:"-"`
	GUID      string     `json:"id"`
	StockLog  int64      `json:"stock_log"`
	StockType string     `json:"stock_type"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *string    `json:"updated_by"`
}
