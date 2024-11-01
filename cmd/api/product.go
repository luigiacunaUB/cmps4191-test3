package main

import (
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
