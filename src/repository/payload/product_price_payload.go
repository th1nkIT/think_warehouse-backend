package payload

import "time"

type readProductPrice struct {
	ID           int64      `json:"-"`
	GUID         string     `json:"id"`
	Price        int64      `json:"price"`
	DiscountType *string    `json:"discount_type"`
	Discount     *int64     `json:"discount"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    string     `json:"created_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
	UpdatedBy    *string    `json:"updated_by"`
}
