package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {

	//setup a new router
	router := httprouter.New()

	//routes
	router.HandlerFunc(http.MethodGet, "/", a.Index)

	return router
}
