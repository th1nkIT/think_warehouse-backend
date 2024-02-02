package payload

import sqlc "think_warehouse/src/repository/pgbo_sqlc"

type GetStockByProductAndVariantParams struct {
	ProductID           string `json:"product_id"`
	ProductVariantID    string `json:"product_variant_id"`
	SetProductVariantID bool   `json:"set_product_variant_id"`
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
