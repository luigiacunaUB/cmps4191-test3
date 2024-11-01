package main

import (
	"fmt"
	"net/http"

	"github.com/luigiacunaUB/cmps4191-test3/internal/data"
	"github.com/luigiacunaUB/cmps4191-test3/internal/validator"
)

func (a *applicationDependencies) createProduct(w http.ResponseWriter, r *http.Request) {
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

	v := validator.New()

	data.ValidateProduct(v, product)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData)
}
