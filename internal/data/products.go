package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	//"github.com/luigiacunaUB/cmps4191-test3/internal/errors"
	"github.com/luigiacunaUB/cmps4191-test3/internal/validator"
)

type ProductModel struct {
	DB     *sql.DB
	logger *slog.Logger
}

type Product struct {
	ID        int64     `json:"productid"`
	ProdName  string    `json:"productname"`
	Category  string    `json:"category"`
	ImgURL    string    `json:"imageurl"`
	AvgRating int       `json:"averagerating"`
	AddedDate time.Time `json:"-"`
}

type ProductRatingModel struct {
	ID     int64 `json:"ratingid"`
	ProdID int64 `json:"prodid"`
	Rating int64 `json:"rating"`
}

func ValidateProduct(v *validator.Validator, p ProductModel, product *Product) {
	v.Check(product.ProdName != "", "productname", "must be provided")
	v.Check(len(product.ProdName) <= 25, "productname", "must not be more than 25 bytes long")

	v.Check(product.Category != "", "category", "must be provided")
	v.Check(len(product.Category) <= 25, "category", "must not be more than 25 bytes long")

	v.Check(product.ImgURL != "", "imageurl", "must be provided")
	v.Check(len(product.ImgURL) <= 100, "imageurl", "must not be more than 100 bytes long")

	exists, err := p.CheckIfProdExist(product.ProdName)

	if err != nil {
		v.Check(false, product.ProdName, "not found")
	}

	if exists {
		v.Check(false, product.ProdName, "found")
	}

}

func (p ProductModel) Insert(product *Product) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside Insert Func")
	logger.Info("Insert parameters", "productname", product.ProdName)
	/*query := `INSERT INTO product(prodname,category,imgurl)
	VALUES ($1,$2,$3,$4)
	RETURNING id,addeddate
	`*/

	query:=`INSERT INTO product(prodname,category,imgurl)
			VALUES('phone','device','theurl');
			INSERT INTO prodratings(prodid,rating) VALUES(1,5);`
	args := []any{product.ProdName, product.Category, product.ImgURL}
	logger.Info("Insert parameters", "productname", args)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	logger.Info("reaches here")
	if p.DB == nil {
		logger.Error("Database connection is nil")
		return fmt.Errorf("database connection is nil")
	}
	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID, &product.AddedDate)
}

func (p ProductModel) Get(id int64) (*Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id,prodname,category,imgurl,category,addeddate
			FROM product
			WHERE id = $1
			`
	var product Product

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.ProdName, &product.AddedDate)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &product, nil
}

func (p ProductModel) Update(product *Product) error {
	query := `UPDATE product
			SET prodname = $1, category = $2, imgurl = $3, rating $4
			WHERE id = $5
			RETURNING id
			`
	args := []any{product.ProdName, product.Category, product.ImgURL, product.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID)
}

func (p ProductModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM product WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (p ProductModel) DisplayAll() ([]Product, error) {
	query := `SELECT id,prodname,category,imgurl,category,addeddate FROM product`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.ProdName, &product.Category, &product.ImgURL, &product.AddedDate)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p ProductModel) CheckIfProdExist(prodname string) (bool, error) {
	query := `SELECT 1 FROM product WHERE prodname = $1 LIMIT 1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := p.DB.QueryRowContext(ctx, query, prodname)
	var exist int
	err := row.Scan(&exist)
	if err == sql.ErrNoRows {
		return false, nil
	} else {
		return true, nil
	}

}
