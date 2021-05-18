package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	fh "github.com/udacity/udagram-restapi-golang/controllers/v0/feed"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	/*
		Create a "Subrouter" dedicated to /api which will use the PathPrefix
		more on nesting routes with gorrila mux here https://stackoverflow.com/questions/25107763/nested-gorilla-mux-router-does-not-work
		https://binx.io/blog/2018/11/27/go-gorilla/
	*/
	apiRouter := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)

	// This step is where we connect our "root" router and our "Subrouter" together.
	r.PathPrefix("/api").Handler(apiRouter)

	// Define "root" routes using r
	r.Methods("GET").Path("/").HandlerFunc(index)

	// Define "Subrouter" routes using apiRouter, prefix is /api
	apiRouter.Methods("GET").Path("/").HandlerFunc(api)
	apiRouter.Methods("GET").Path("/v0/feed").HandlerFunc(fh.IndexHandler)
	apiRouter.Methods("POST").Path("/v0/feed").HandlerFunc(fh.CreateFeedItemHandler)
	apiRouter.Methods("GET").Path("/v0/feed/{id}").HandlerFunc(fh.GetFeedItemHandler)
	apiRouter.Methods("GET").Path("/v0/feed/signed-url/{fileName}").HandlerFunc(fh.GetGetSignedUrlHandler)

	http.ListenAndServe(":"+port, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/api/v0/")
}

func api(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "v0")
}
