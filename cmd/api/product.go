package main

import (
	"fmt"
	"net/http"
)

func (a *applicationDependencies) createProduct(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		ProductName string `json:"productname`
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData)
}
