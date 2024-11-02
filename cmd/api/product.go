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
	var incomingData struct {
		ProdName string `json:"productname"`
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	product := &data.Product{
		ProdName: incomingData.ProdName,
	}
	a.logger.Info(incomingData.ProdName)

	v := validator.New()

	data.ValidateProduct(v, product)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	} else {
		a.logger.Info("Validate Product Pass")
	}

	err = a.ProductModel.Insert(product)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	} else {
		a.logger.Info("Insert Pass")
	}
	fmt.Fprintf(w, "%+v\n", incomingData)
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/product/%d", product.ID))

	data := envelope{
		"productname": product,
	}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayProductHandler(w http.ResponseWriter, r *http.Request) {
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
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

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
	}

	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if incomingData.ProdName != nil {
		product.ProdName = *incomingData.ProdName
	}

	v := validator.New()
	data.ValidateProduct(v, product)
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

func (a *applicationDependencies)deleteProductHandler(w http.ResponseWriter,r *http.Request){
	id,err := a.readIDParam(r)
	if err != nil{
		a.notFoundResponse(w,r)
		return
	}

	err = a.ProductModel.Delete(id)

	if err != nil{
		switch{
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w,r)
		default:
			a.serverErrorResponse(w,r,err)
		}
	}
	data := envelope{
		"message":"product deleted",
	}

	err = a.writeJSON(w,http.StatusOK,data,nil)
	if err != nil{
		a.serverErrorResponse(w,r,err)
	}
}
