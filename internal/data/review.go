/*
review table workflow
 1. create a review for a specfic product: route /review/"reviewid"/"prodid"
    a.Check if a product exist, if not stop the process
    b.If product is found add the review with no rating
 2. display a review for a specfic product:
    a. route
*/
package data

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	//"github.com/luigiacunaUB/cmps4191-test3/internal/errors"
	"github.com/luigiacunaUB/cmps4191-test3/internal/validator"
)

type ReviewModel struct {
	DB *sql.DB
}

type Review struct {
	ID              int    `json:"id"`
	ProdDid         int    `json:"productid"`
	ProductName     string `json:"productname"`
	Review          string `json:"review"`
	HelpfullCounter int    `json:"helpful"`
	AddedDate       string `json:"-"`
}

func ValidateReview(v *validator.Validator, r ReviewModel, review *Review) {
	//check if the review is not empty
	v.Check(review.Review != "", "review", "must be provided")
	v.Check(len(review.Review) <= 100, "review", "must not be more than 100 bytes long")
	//check if the helpful counter is between 1 and 5
	v.Check(review.HelpfullCounter >= 1 && review.HelpfullCounter <= 5, "rating", "must be between 1 and 5")

}

func (r ReviewModel) InsertReview(review *Review, instruction bool) error {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside InsertReview Func")
	//-----------------------------------------------------------------------------
	query := `INSERT INTO review(prodid,review) VALUES($1,$2)`
	args := []any{review.ProdDid, review.Review}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return r.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID, &review.ProductName)

}

func (r ReviewModel) CheckIfProdIDExist(prodid int) (bool, error) {
	query := `SELECT 1 FROM product WHERE id = $1 LIMIT 1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := r.DB.QueryRowContext(ctx, query, prodid)
	var exist int
	err := row.Scan(&exist)
	if err == sql.ErrNoRows {
		return false, nil
	} else {
		return true, nil
	}

}
