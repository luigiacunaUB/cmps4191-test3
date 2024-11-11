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

// deprecated
func ValidateRating(v *validator.Validator, r ReviewModel, review *Review) {
	//check if the helpful counter is between 1 and 5
	v.Check(review.HelpfullCounter >= 1 && review.HelpfullCounter <= 5, "rating", "must be between 1 and 5")

}

func ValidateHelpfulAnswer(v *validator.Validator, answer string) {
	v.Check(answer == "yes" || answer == "Yes" || answer == "YES" || answer == "y" || answer == "Y", "helpful", "must be 'yes' or 'y'")
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

func (r ReviewModel) CheckIfReviewExist(reviewID int) (bool, error) {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside CheckIfReviewExist")
	logger.Info("Review ID: ", reviewID)
	query := `SELECT id FROM review WHERE id = $1 LIMIT 1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	logger.Info("Reaches Here One")
	var exists bool
	logger.Info("Reaches Here Two")
	err := r.DB.QueryRowContext(ctx, query, reviewID).Scan(&exists)
	logger.Info("Reaches Here Three")
	logger.Info("Exist status: ", exists)
	if err == nil && exists == true {
		logger.Info("Review ID FOUND")
		return true, nil
	} else if err == nil && exists == false {
		logger.Info("Review ID NOT FOUND")
		return false, nil
	}
	logger.Info("Exiting")
	return true, nil
}
func (r ReviewModel) DisplaySpecificReview(productID int, ReviewID int) (*Review, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside DisplaySpecficReview Func")
	tx, err := r.DB.Begin()
	if err != nil {
		logger.Error("Cannot begin Transaction")
		return nil, err
	}
	review := &Review{}
	//first query
	firstquery := `SELECT id,prodname FROM product WHERE id=$1`
	err = tx.QueryRow(firstquery, productID).Scan(&review.ProdDid, &review.ProductName)
	if err != nil {
		logger.Error("SELECT product FAILED ABORTING!")
		return nil, err
	}
	//second query
	secondquery := `SELECT id,review FROM review WHERE id=$1`
	err = tx.QueryRow(secondquery, ReviewID).Scan(&review.ID, &review.Review)
	err = tx.Commit()
	if err != nil {
		logger.Info("Error getting review info")
		return nil, err
	}
	return review, nil

}

func (r ReviewModel) UpdateSpecificReview(ReviewID int, data string) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside UpdateSpecficReview Func")

	logger.Info("Review ID inside update sql", ReviewID)
	logger.Info("data to update: ", data)

	query := `UPDATE review SET review = $1 WHERE id = $2`
	args := []any{data, ReviewID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
func (r ReviewModel) DeleteReview(ReviewID int) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside DeleteSpecficReview Func")
	query := `DELETE FROM review WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := r.DB.ExecContext(ctx, query, ReviewID)
	if err != nil {
		return err
	}

	return nil
}

func (r ReviewModel) DisplayAllReviews(review string) ([]Review, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside DisplayAll Func")
	var reviews []Review
	query := `SELECT
    	p.id AS product_id,
    	p.prodname AS product_name,
    	r.id AS review_id,
    	r.review,
		COALESCE(r.helpfulcounter,0) AS helpfulcounter
	FROM
    	product p
	JOIN
    	review r ON p.id = r.prodid
	WHERE
    to_tsvector('english', r.review) @@ plainto_tsquery('english', $1)
    OR $1 = ''`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, review)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var review Review
		err := rows.Scan(&review.ProdDid, &review.ProductName, &review.ID, &review.Review, &review.HelpfullCounter)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil

}

func (r ReviewModel) DisplayAllReviewsForProduct(productID int) ([]Review, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside DisplayAllReviewsForProduct Func")

	var reviews []Review
	query := `SELECT
                p.id AS product_id,
                p.prodname AS product_name,
                r.id AS review_id,
                r.review
              FROM
                product p
              JOIN
                review r ON p.id = r.prodid
              WHERE
                p.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var review Review
		err := rows.Scan(&review.ProdDid, &review.ProductName, &review.ID, &review.Review)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r ReviewModel) HelpfulAnswerAdd(id int64) bool {
	// Prepare SQL query to increment helpfulcounter for a specific review
	query := `
	 UPDATE review
	 SET helpfulcounter = COALESCE(helpfulcounter, 0) + 1
	 WHERE id = $1
 `

	// Execute the query
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return false
	}

	return true
}
