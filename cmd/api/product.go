package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *applicationDependencies) createProduct(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		ProductName string `json:"productname`
	}

	err := json.NewDecoder(r.Body).Decode(&incomingData)
	if err != nil {
		a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData)
}
