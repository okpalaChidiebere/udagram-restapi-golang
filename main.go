package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	fh "github.com/udacity/udagram-restapi-golang/controllers/v0/feed"
	uh "github.com/udacity/udagram-restapi-golang/controllers/v0/users"
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
	indexRouter := r.PathPrefix("/api/v0").Subrouter().StrictSlash(true)

	// This step is where we connect our "index" SubRouter to Feed SubRouter and Users SubRouter
	feedRouter := indexRouter.PathPrefix("/feed").Subrouter().StrictSlash(true)
	usersRouter := indexRouter.PathPrefix("/users").Subrouter().StrictSlash(true)

	authRouter := usersRouter.PathPrefix("/auth").Subrouter()

	// Define "Subrouter" routes using indexRouter
	indexRouter.Methods("GET").Path("/").HandlerFunc(indexRouterHandler)

	// Define "root" routes using r
	r.Methods("GET").Path("/").HandlerFunc(index)

	// Define "Subrouter" routes using feedRouter, prefix is /api/v0/feed/...
	feedRouter.Methods("GET").Path("/").HandlerFunc(fh.IndexHandler)
	feedRouter.Methods("POST").Path("/").HandlerFunc(fh.CreateFeedItemHandler)
	feedRouter.Methods("GET").Path("/{id}").HandlerFunc(fh.GetFeedItemHandler)
	feedRouter.Methods("GET").Path("/signed-url/{fileName}").HandlerFunc(fh.GetGetSignedUrlHandler)

	// Define "Subrouter" routes using usersRouter, prefix is /api/v0/users/...
	usersRouter.Methods("GET").Path("/{id}").HandlerFunc(uh.GetUserHandler)

	authRouter.Methods("GET").Path("").HandlerFunc(authIndex)               // ../api/v0/users/auth it was a little bit odd that i did not specify the index route "/" in the Path() method. But it works :)
	authRouter.Methods("POST").Path("").HandlerFunc(uh.RegisterUserHandler) // ../api/v0/users/auth

	http.ListenAndServe(":"+port, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/api/v0/")
}

func indexRouterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "v0")
}

func authIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "auth")
}
