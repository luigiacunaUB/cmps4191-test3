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

}

func ValidateRating(v *validator.Validator, r ReviewModel, review *Review) {
	//check if the helpful counter is between 1 and 5
	v.Check(review.HelpfullCounter >= 1 && review.HelpfullCounter <= 5, "rating", "must be between 1 and 5")

}

func (r ReviewModel) InsertReview(review *Review, instruction bool) error {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside InsertReview Func")
	//-----------------------------------------------------------------------------
	//double query to get correct into
	//1. Do the actual insert, then grab the product name
	tx, err := r.DB.Begin()
	if err != nil {
		logger.Error("Cannot begin Transaction")
		return err
	}
	//first query
	firstquery := `INSERT INTO review(prodid,review) VALUES($1,$2) RETURNING id`
	err = tx.QueryRow(firstquery, review.ProdDid, review.Review).Scan(&review.ID)
	if err != nil {
		logger.Error("INSERT INTO review FAILED ABORTING!")
		return err
	}
	secondquery := `SELECT prodname FROM product WHERE id=$1`
	err = tx.QueryRow(secondquery, review.ProdDid).Scan(&review.ProductName)
	if err != nil {
		logger.Error("SELECT product FAILED ABORTING!")
		return err
	}
	err = tx.Commit()
	if err != nil {
		logger.Error("INSERT INTO prodrating FAILED ABORTING!")
		//need to implement delete if the insert did pass to keep db clean
		return err
	}
	return nil

}

func (r ReviewModel) CheckIfProdIDExist(prodid int) (bool, error) {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside CheckIfProdIDExist")
	logger.Info("Product ID: ", prodid)
	//query := `SELECT 1 FROM product WHERE id = $1 LIMIT 1`
	query := `SELECT id FROM product WHERE id = $1 LIMIT 1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	logger.Info("Reaches Here One")
	var exists bool
	logger.Info("Reaches Here Two")
	err := r.DB.QueryRowContext(ctx, query, prodid).Scan(&exists)
	logger.Info("Reaches Here Three")
	logger.Info("Exist status: ", exists)
	if err == nil && exists == true {
		logger.Info("Product ID FOUND")
		return true, nil
	} else if err == nil && exists == false {
		logger.Info("Product ID NOT FOUND")
		return false, nil
	}
	logger.Info("Exiting")
	return true, nil
}
