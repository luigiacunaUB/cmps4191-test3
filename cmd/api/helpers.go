package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (a *applicationDependencies) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	jsResponse, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	jsResponse = append(jsResponse, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	_, err = w.Write(jsResponse)

	if err != nil {
		return err
	}
	return nil
}

func (a *applicationDependencies) readJSON(w http.ResponseWriter, r *http.Request, destination any) error {

	maxBytes := 256_000
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode((destination))
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var InvalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("the body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("the body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("the body contains the incorrent JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("the body contains the incorrect JSON type (at charater %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("The body must not be empty")
		case strings.HasPrefix(err.Error(), "json:unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json:unknown field")
			return fmt.Errorf("the body contains unknown key  %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("the body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &InvalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("the body must only contain a single JSON value")
	}
	return nil
}

func (a *applicationDependencies) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (a *applicationDependencies) prodAndIDChecks(w http.ResponseWriter, r *http.Request) (bool, int, int) {
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
	if err != nil && exist == false {
		a.serverErrorResponse(w, r, err) //need to update error to say product does not exist
		return false, 0, 0
	} else if err == nil && exist == true {
		a.logger.Info("Product ID Pass")
		return true, productID, reviewID
	}
	exist, err = a.ReviewModel.CheckIfReviewExist(reviewID)
	if err != nil && exist == false {
		a.serverErrorResponse(w, r, err) //need to update error to say review does not exist
		return false, 0, 0
	} else if err == nil && exist == true {
		a.logger.Info("Review ID Pass")
		return true, productID, reviewID
	}
	return false, 0, 0
}
