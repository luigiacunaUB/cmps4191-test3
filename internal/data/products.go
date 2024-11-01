package data

import (
	"time"

	"github.com/luigiacunaUB/cmps4191-test3/internal/validator"
)

type Product struct {
	ID        int64     `json:"id"`
	ProdName  string    `json:"productname"`
	AddedDate time.Time `json:"-"`
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.ProdName != "", "productname", "must be provided")
	v.Check(len(product.ProdName) <= 25, "productname", "must not be more than 25 bytes long")
}
