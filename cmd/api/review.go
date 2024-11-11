package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func (a *applicationDependencies) displaySpecficReviewHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info(r.URL.String())
	var pID, rID string
	parts := strings.Split(r.URL.Path, "/")
	//spit the url string to only grab the values needed
	if len(parts) >= 5 {
		pID = parts[2] // product ID is at index 2 in the path
		rID = parts[4] // review ID is at index 4 in the path

		fmt.Fprintf(w, "Product ID: %s\n", pID)
		fmt.Fprintf(w, "Review ID: %s\n", rID)
	}
	var productID, reviewID int
	productID, err := strconv.Atoi(pID)
	reviewID, err = strconv.Atoi(rID)

	exist, err := a.ReviewModel.CheckIfProdIDExist(productID)
	if err == nil && exist == false {
		a.serverErrorResponse(w, r, err) //need to update error to say product does not exist
		return
	} else if err == nil && exist == true {
		a.logger.Info("Product ID Pass")
	}
	exist, err = a.ReviewModel.CheckIfReviewExist(reviewID)
	if err == nil && exist == false {
		a.serverErrorResponse(w, r, err) //need to update error to say review does not exist
		return
	} else if err == nil && exist == true {
		a.logger.Info("Review ID Pass")
	}

	review, err := a.ReviewModel.DisplaySpecificReview(productID, reviewID)
	data := envelope{
		"product":    review.ProductName,
		"product ID": review.ProdDid,
		"review ID":  review.ID,
		"review":     review.Review,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}
func (a *applicationDependencies) updateSpecficReviewHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Inside update handeler")
	pass, productID, reviewID := a.prodAndIDChecks(w, r)
	if pass == false {
		a.logger.Info("Failed") //need to debug all checks still passing even with errors
		return
	}
	review, err := a.ReviewModel.DisplaySpecificReview(productID, reviewID) //everything from review here should come here
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
	//a.logger.Info("1. Original Review: ", review.Review)

	var incomingData struct {
		Review *string `json:"review"`
	}
	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	//a.logger.Info("2. Modified Review: ", *incomingData.Review)
	if incomingData.Review != nil {
		review.Review = *incomingData.Review
	}
	v := validator.New()
	data.ValidateReview(v, a.ReviewModel, review)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
	}
	//call update
	//need to call productID,reviewID,incomingdata.Review
	//a.logger.Info("3. Original Review: ", review.Review)
	//a.logger.Info("4. Modified Review: ", *incomingData.Review)
	err = a.ReviewModel.UpdateSpecificReview(reviewID, review.Review)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"Review updated for ID": review.ID,
		"Updated Review:":       review.Review,
		"Product ID":            review.ProdDid,
		"Product":               review.ProductName,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}

}

func (a *applicationDependencies) deleteSpecficReviewHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Inside delete handeler")
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}
	var convertedInt int64 = int64(id)
	a.logger.Info("Review id to be deleted: ", convertedInt)
	//call delete
	err = a.ReviewModel.DeleteReview(int(convertedInt))
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}

	data := envelope{
		"review": "deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayAllReviewsHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Inside Display Handler")
	review, err := a.ReviewModel.DisplayAllReviews()
	if err != nil {
		switch {
		case errors.Is(err, data.QueryFail):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)

		}
	}
	reviewCount := len(review)
	a.logger.Info("Review table Count: ", reviewCount)

	if reviewCount == 0 {
		data := envelope{
			"message": "there is no reviews",
		}
		err = a.writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			a.serverErrorResponse(w, r, err)
		}
	} else {
		data := envelope{
			"Review": review,
		}
		err = a.writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			a.serverErrorResponse(w, r, err)
		}
	}
}

func (a *applicationDependencies) displayAllReviewsForSpecificProductHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Inside Display All Reviews for a specfic Product")
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}
	var convertedInt int64 = int64(id)
	review, err := a.ReviewModel.DisplayAllReviewsForProduct(int(convertedInt))
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	reviewCount := len(review)
	a.logger.Info("Review table Count: ", reviewCount)

	if reviewCount == 0 {
		data := envelope{
			"message": "there is no reviews",
		}
		err = a.writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			a.serverErrorResponse(w, r, err)
		}
	} else {
		data := envelope{
			"Review": review,
		}
		err = a.writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			a.serverErrorResponse(w, r, err)
		}
	}

}
