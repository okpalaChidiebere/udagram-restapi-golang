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

	//CORS Should be restricted
	r.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8100")
				w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
				w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,PUT,PATCH,POST,DELETE")
				next.ServeHTTP(w, r)
			})
		},
	)

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
	authFeedRouter := feedRouter.PathPrefix("").Subrouter() //this subRouter under feedRouter
	/*
		What we are doing here is we are type casting custom type Adapter to a mux.Middleware type
		the root type for our Adapter is func(http.Handler) http.Handler which is thesame to Middleware type. This is why we are able to do this
		https://pkg.go.dev/github.com/gorilla/mux@v1.8.0#MiddlewareFunc
	*/
	authFeedRouter.Use(mux.MiddlewareFunc(uh.RequireAuthHandler())) //we make the subrouter from feedRouter to be protected

	//Define not protected feedRouter. The RequireAuthHandler middleare will not apply to them
	feedRouter.Methods("GET").Path("").HandlerFunc(fh.IndexHandler)

	//Define protected feedRouter
	authFeedRouter.Methods("POST").Path("").HandlerFunc(fh.CreateFeedItemHandler)
	authFeedRouter.Methods("GET").Path("/{id}").HandlerFunc(fh.GetFeedItemHandler)
	authFeedRouter.Methods("GET").Path("/signed-url/{fileName}").HandlerFunc(fh.GetGetSignedUrlHandler)

	// Define "Subrouter" routes using usersRouter, prefix is /api/v0/users/...
	usersRouter.Methods("GET").Path("/{id}").HandlerFunc(uh.GetUserHandler)

	authRouter.Methods("GET").Path("").HandlerFunc(authIndex)                  // ../api/v0/users/auth it was a little bit odd that i did not specify the index route "/" in the Path() method. But it works :)
	authRouter.Methods("POST").Path("").HandlerFunc(uh.RegisterUserHandler)    // ../api/v0/users/auth
	authRouter.Methods("POST").Path("/login").HandlerFunc(uh.LoginUserHandler) // ../api/v0/users/auth/login
	authRouter.Methods("GET").Path("/verification").HandlerFunc(uh.Adapt(
		http.HandlerFunc(uh.VerificationHandler),
		uh.RequireAuthHandler(),
	).ServeHTTP) // ../api/v0/users/auth/verification

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
