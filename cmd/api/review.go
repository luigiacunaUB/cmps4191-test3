package main

import (
	"fmt"
	"net/http"

	"github.com/luigiacunaUB/cmps4191-test3/internal/data"
	"github.com/luigiacunaUB/cmps4191-test3/internal/validator"
)

func (a *applicationDependencies) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Inside createReview")
	//incoming data that will accept the request from curl command, a creator of review cannot rate itself
	//only the productID and Review will be accepted
	var incomingData struct {
		ProdDid int    `json:"productid"`
		Review  string `json:"review"`
	}
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	//prep the incoming data to validated
	review := &data.Review{
		ProdDid: incomingData.ProdDid,
		Review:  incomingData.Review,
	}

	//call the validator
	v := validator.New()
	//check the incoming data if it passes validation
	data.ValidateReview(v, a.ReviewModel, review)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	} else {
		a.logger.Info("Validation Pass")
	}

	//check if a product exist based on id
	exist, err := a.ReviewModel.CheckIfProdIDExist(incomingData.ProdDid)
	if err == nil && exist == false {
		a.serverErrorResponse(w, r, err) //need to update error to say product does not exist
		return
	} else if err == nil && exist == true {
		a.logger.Info("Product ID Pass")
	}

	err = a.ReviewModel.InsertReview(review, exist)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	fmt.Fprintf(w, "%+v\n", incomingData)
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/product/%d/review/%d", review.ProdDid, review.ID))

	data := envelope{
		"ReviewID":     review.ID,
		"ProductID":    review.ProdDid,
		"Product Name": review.ProductName,
		"Review":       review.Review,
		"Location":     fmt.Sprintf("localhost:4000/product/%d/review/%d", review.ProdDid, review.ID),
	}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
