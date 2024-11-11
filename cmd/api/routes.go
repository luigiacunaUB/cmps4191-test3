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
	//basic
	router.HandlerFunc(http.MethodGet, "/", a.Index)                         //root page
	router.HandlerFunc(http.MethodGet, "/healthcheck", a.healthCheckHandler) //healthcheck
	//routes dealing with product----------------------------------------------------------------------------
	router.HandlerFunc(http.MethodPost, "/product", a.createProduct)                //add product
	router.HandlerFunc(http.MethodGet, "/product/:id", a.displayProductHandler)     //display product
	router.HandlerFunc(http.MethodPatch, "/product/:id", a.updateProductHandler)    //update product
	router.HandlerFunc(http.MethodDelete, "/product/:id", a.deleteProductHandler)   //delete product
	router.HandlerFunc(http.MethodGet, "/products/all", a.displayAllProductHandler) //display all products
	//routes dealing with reviews----------------------------------------------------------------------------
	router.HandlerFunc(http.MethodPost, "/addreview/products/:id", a.createReviewHandler)
	router.HandlerFunc(http.MethodGet, "/product/:id/review/:id", a.displaySpecficReviewHandler)
	router.HandlerFunc(http.MethodPost, "/product/:id/review/:id", a.updateSpecficReviewHandler)
	router.HandlerFunc(http.MethodDelete, "/review/delete/:id", a.deleteSpecficReviewHandler)
	router.HandlerFunc(http.MethodGet, "/allreviews", a.displayAllReviewsHandler) //this one
	router.HandlerFunc(http.MethodGet, "/allreviews/product/:id", a.displayAllReviewsForSpecificProductHandler)
	router.HandlerFunc(http.MethodPatch, "/inserthelpful/:id", a.incrementHelpfulCounter)

	return a.recoverPanic(router)
}
