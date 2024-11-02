package data

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/luigiacunaUB/cmps4191-test3/internal/validator"
)

type ProductModel struct {
	DB     *sql.DB
	logger *slog.Logger
}

type Product struct {
	ID        int64     `json:"id"`
	ProdName  string    `json:"productname"`
	AddedDate time.Time `json:"-"`
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.ProdName != "", "productname", "must be provided")
	v.Check(len(product.ProdName) <= 25, "productname", "must not be more than 25 bytes long")
}

func (p ProductModel) Insert(product *Product) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside Insert Func")
	logger.Info("Insert parameters", "productname", product.ProdName)
	query := `INSERT INTO product(prodname)
	VALUES ($1)
	RETURNING id,addeddate
	`

	args := []any{product.ProdName}
	logger.Info("Insert parameters", "productname", args)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	logger.Info("reaches here")
	if p.DB == nil {
		logger.Error("Database connection is nil")
		return fmt.Errorf("database connection is nil")
	}
	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID,&product.AddedDate)
}
