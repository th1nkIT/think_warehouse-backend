package payload

import "time"

type readProductStock struct {
	ID        int64      `json:"-"`
	GUID      string     `json:"id"`
	Stock     int64      `json:"stock"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *string    `json:"updated_by"`
}
