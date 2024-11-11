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

type ProductModel struct {
	DB     *sql.DB
	logger *slog.Logger
}

// The Product Model now contains fields of both product and prodratings table, setting
type Product struct {
	ID        int64     `json:"productid"`
	ProdName  string    `json:"productname"`
	Category  string    `json:"category"`
	ImgURL    string    `json:"imageurl"`
	Rating    int       `json:"rating"`
	AddedDate time.Time `json:"-"`
}

func ValidateProduct(v *validator.Validator, p ProductModel, product *Product) {

	v.Check(product.ProdName != "", "productname", "must be provided")
	v.Check(len(product.ProdName) <= 25, "productname", "must not be more than 25 bytes long")

	v.Check(product.Category != "", "category", "must be provided")
	v.Check(len(product.Category) <= 25, "category", "must not be more than 25 bytes long")

	v.Check(product.ImgURL != "", "imageurl", "must be provided")
	v.Check(len(product.ImgURL) <= 100, "imageurl", "must not be more than 100 bytes long")

	v.Check(product.Rating >= 1 && product.Rating <= 5, "rating", "must be between 1 and 5")

}

func (p ProductModel) Insert(product *Product, instruction bool) error {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside Insert Func")
	//---------------------------------------------------------------------------------
	//set the senario up do insert of product and rating
	if instruction == false {
		//using a transaction to process two queries, if one fails rollback
		logger.Info("Inserting a new product and adding a rating")
		tx, err := p.DB.Begin()
		if err != nil {
			logger.Error("Cannot begin Transaction")
			return err
		}

		//first query
		firstquery := `INSERT INTO product(prodname,category,imgurl) VALUES($1,$2,$3) RETURNING id,addeddate`
		err = tx.QueryRow(firstquery, product.ProdName, product.Category, product.ImgURL).Scan(&product.ID, &product.AddedDate)
		//incase it errors out
		if err != nil {
			logger.Error("INSERT INTO product FAILED ABORTING!")
			return err
		}

		secondquery := `INSERT INTO prodratings(prodid,rating) VALUES($1,$2)`
		_, err = tx.Exec(secondquery, product.ID, product.Rating)
		err = tx.Commit()
		if err != nil {
			logger.Error("INSERT INTO prodrating FAILED ABORTING!")
			//need to implement delete if the insert did pass to keep db clean
			return err
		}
		//------------------------------------------------------------------------------------------
	} else if instruction == true {
		//only do the rating
		logger.Info("Only adding a rating for a existing product")
		//need to do a double query
		//1. Get the item id for the existing id and pass into a variable
		//2. with the variable created make the insertion

		tx, err := p.DB.Begin()
		if err != nil {
			logger.Error("Cannot begin Transaction")
			return err
		}

		//first query getting the product id in question to pass into a variable
		firstquery := `SELECT id FROM product WHERE  prodname=($1)`
		err = tx.QueryRow(firstquery, product.ProdName).Scan(&product.ID)
		//incase it errors out
		if err != nil {
			logger.Error("SELECT failed")
			return err
		}
		secondquery := `INSERT INTO prodratings(prodid,rating) VALUES($1,$2)`
		_, err = tx.Exec(secondquery, product.ID, product.Rating)
		err = tx.Commit()
		if err != nil {
			logger.Error("INSERT INTO prodrating FAILED ABORTING!")
			//need to implement delete if the insert did pass to keep db clean
			return err
		}

	}
	return nil

}

// -----------------------------------------------------------------------------------------------------------------------------
func (p ProductModel) Get(id int64) (*Product, error) {
	//GET: Gets the item info for a specfic item along with its average rating, coming from prodratings table
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//----------------------------------------------------------------------------------------------
	//first get the product info from the product table
	//using transaction to excute mutiple quries
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	tx, err := p.DB.Begin()
	if err != nil {
		logger.Info("Cannot Begin Transaction")
	}
	var product Product //var hold info needed
	//first query get the item info
	selectQuery := `SELECT id,prodname,category,imgurl,addeddate FROM product WHERE id = $1`
	err = tx.QueryRow(selectQuery, id).Scan(&product.ID, &product.ProdName, &product.Category, &product.ImgURL, &product.AddedDate)
	if err != nil {
		logger.Error("ERROR getting product info")
	}
	//---------------------------------------------------------------------------------------------------------
	//second query to pull the average rating
	averageQuery := `SELECT COALESCE(ROUND(AVG(rating)), 0) AS average_rating FROM prodratings WHERE prodid = $1`

	err = tx.QueryRow(averageQuery, id).Scan(&product.Rating)
	err = tx.Commit()
	if err != nil {
		logger.Error("ERROR getting average rating")
	}
	logger.Info("Info", product)

	return &product, nil
}

// ------------------------------------------------------------------------------------------------------------------------------------
func (p ProductModel) Update(product *Product) error {
	query := `UPDATE product
			SET prodname = $1, category = $2, imgurl = $3
			WHERE id = $4
			RETURNING id
			`
	args := []any{product.ProdName, product.Category, product.ImgURL, product.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID)
}

// --------------------------------------------------------------------------------------------------------------------------------------
func (p ProductModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	tx, err := p.DB.Begin()
	if err != nil {
		logger.Info("Cannot Begin Transaction")
	}
	//for the delete function it will go in reverse, deleting all sub table info, before main table infor
	//deleting from prodratings
	DelProdRatingsInfor := `DELETE FROM prodratings WHERE prodid = $1`
	ratingresult, err := tx.Exec(DelProdRatingsInfor, id)
	if err != nil {
		logger.Info("Cannot delete from prodratings")
		tx.Rollback()
	}
	DelProd := `DELETE FROM product WHERE id = $1`
	productResult, err := tx.Exec(DelProd, id)
	if err != nil {
		logger.Info("Cannot delete from product")
		tx.Rollback()
	}
	err = tx.Commit()
	if err != nil {
		logger.Info("Error DELETING")
	}

	ratingsrAffected, err := ratingresult.RowsAffected()
	if err != nil {
		return err
	}

	rowsAffected, err := productResult.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 || ratingsrAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}

// ------------------------------------------------------------------------------------------------------------------------------
func (p ProductModel) DisplayAll(productname string, category string) ([]Product, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Inside DisplayAll")
	var products []Product
	//---------------------------------------------------------------------------------

	// Query to get all product data and the overall average rating for all products
	//orginal query
	/*query := `
	  SELECT
	      p.id, p.prodname, p.category, p.imgurl, p.addeddate,
	      ROUND((SELECT AVG(r.rating) FROM prodratings r WHERE r.prodid = p.id)) AS rating
	  FROM
	      product p
	  `*/

	query := `SELECT 
    p.id, 
    p.prodname, 
    p.category, 
    p.imgurl, 
    p.addeddate, 
    ROUND((SELECT AVG(r.rating) FROM prodratings r WHERE r.prodid = p.id)) AS rating
FROM 
    product p
WHERE 
    (to_tsvector('english', p.prodname) @@ plainto_tsquery('english', COALESCE($1, ''))
    OR $1 = '')
    AND 
    (to_tsvector('english', p.category) @@ plainto_tsquery('english', COALESCE($2, ''))
    OR $2 = '');
`

	// Set up a timeout context for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	rows, err := p.DB.QueryContext(ctx, query, productname, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows and scan the results into the products slice
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.ProdName, &product.Category, &product.ImgURL, &product.AddedDate, &product.Rating)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return the products slice
	return products, nil

	/*query := `SELECT id,prodname,category,imgurl,category,addeddate FROM product`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	return products, nil*/
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
