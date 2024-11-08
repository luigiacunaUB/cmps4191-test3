package main

import (
	"fmt"
	"net/http"

	"github.com/luigiacunaUB/cmps4191-test3/internal/data"
	"github.com/luigiacunaUB/cmps4191-test3/internal/validator"
)

func (a *applicationDependencies) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Inside createReview")
	var incomingData struct {
		ProdDid         int    `json:"productid"`
		Review          string `json:"review"`
		HelpfullCounter int    `json:"helpful"`
	}
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	review := &data.Review{
		ProdDid:         incomingData.ProdDid,
		Review:          incomingData.Review,
		HelpfullCounter: incomingData.HelpfullCounter,
	}
	//a.logger.Info(incomingData.ProdDid, incomingData.Review, incomingData.HelpfullCounter)

	v := validator.New()

	data.ValidateReview(v, a.ReviewModel, review)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	} else {
		a.logger.Info("Validation Pass")
	}
	//check if a product exist based on id
	exist, err := a.ReviewModel.CheckIfProdIDExist(incomingData.ProdDid)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
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
		"reviewID":        review.ID,
		"ProductID":       review.ProdDid,
		"Product Name":    review.ProductName,
		"Review":          review.Review,
		"Helpful Counter": review.HelpfullCounter,
	}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
