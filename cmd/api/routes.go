package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {

	//setup a new router
	router := httprouter.New()

	//errors
	//404
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	//405
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)
	//routes
	router.HandlerFunc(http.MethodGet, "/", a.Index)                         //root page
	router.HandlerFunc(http.MethodGet, "/healthcheck", a.healthCheckHandler) //healthcheck
	router.HandlerFunc(http.MethodPost, "/product", a.createProduct)		 //add product
	router.HandlerFunc(http.MethodGet,"/product/:id",a.displayProductHandler)

	return a.recoverPanic(router)
}
