package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/luigiacunaUB/cmps4191-test3/internal/data"
	"github.com/luigiacunaUB/cmps4191-test3/internal/validator"
)

func (a *applicationDependencies) createProduct(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Inside createProduct")
	var incomingData struct { //pushes into
		ProdName string `json:"productname"` //product table
		Category string `json:"category"`    //product table
		ImgURL   string `json:"imageurl"`    //product table
		Rating   int    `json:"rating"`      //prodratings table
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	product := &data.Product{
		ProdName: incomingData.ProdName,
		Category: incomingData.Category,
		ImgURL:   incomingData.ImgURL,
		Rating:   incomingData.Rating,
	}
	a.logger.Info(incomingData.ProdName, incomingData.Category, incomingData.ImgURL, incomingData.Category)

	v := validator.New()

	data.ValidateProduct(v, a.ProductModel, product)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	} else {
		a.logger.Info("Validate Product Pass")
	}

	// Check if the product exists
	exist, err := a.ProductModel.CheckIfProdExist(incomingData.ProdName)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Insert product depending on existence status
	err = a.ProductModel.Insert(product, exist)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData)
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/product/%d", product.ID))

	data := envelope{
		"productname": product.ProdName,
		"category":    product.Category,
		"rating":      product.Rating,
		"id":          product.ID,
	}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayProductHandler(w http.ResponseWriter, r *http.Request) {
	//a.logger.Info("Inside displayProductHandler")

	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	product, err := a.ProductModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
	}

	data := envelope{
		"productname": product,
	}
	//a.logger.Info("pass here")
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}

// ----------------------------------------------------------------------------------------------------------
func (a *applicationDependencies) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}
	product, err := a.ProductModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
	}

	var incomingData struct {
		ProdName *string `json:"productname"`
		Category *string `json:"category"`
		ImgURL   *string `json:"imageurl"`
	}

	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if incomingData.ProdName != nil {
		product.ProdName = *incomingData.ProdName
		product.Category = *incomingData.Category
		product.ImgURL = *incomingData.ImgURL
	}

	v := validator.New()
	data.ValidateProduct(v, a.ProductModel, product)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.ProductModel.Update(product)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"productname": product,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}

}

func (a *applicationDependencies) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.ProductModel.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
	}
	data := envelope{
		"message": "product deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayAllProductHandler(w http.ResponseWriter, r *http.Request) {
	var queryParametersData struct {
		ProdName string
		Category string
		data.Filters
	}

	queryParameters := r.URL.Query()

	queryParametersData.ProdName = a.getSingleQueryParameter(queryParameters, "productname", "")
	queryParametersData.Category = a.getSingleQueryParameter(queryParameters, "category", "")

	v := validator.New()

	queryParametersData.Filters.Page = a.getSingleIntegerParameter(queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameter(queryParameters, "page_size", 10, v)

	queryParametersData.Filters.Sort = a.getSingleQueryParameter(queryParameters, "sort", "id")
	queryParametersData.Filters.SortSafeList = []string{"id", "category", "-id", "-category"}

	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	//products, err := a.ProductModel.DisplayAll()
	products, metadata, err := a.ProductModel.DisplayAll(queryParametersData.ProdName, queryParametersData.Category, queryParametersData.Filters)

	if err != nil {
		switch {
		case errors.Is(err, data.QueryFail):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)

		}
	}
	prodCount := len(products)
	a.logger.Info("Product Table Count: ", prodCount)

	if prodCount == 0 {
		data := envelope{
			"message": "there is no products",
		}
		err = a.writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			a.serverErrorResponse(w, r, err)
		}
	} else {
		data := envelope{
			"Product":   products,
			"@metadata": metadata,
		}
		err = a.writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			a.serverErrorResponse(w, r, err)
		}
	}

}
