package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	fh "github.com/udacity/udagram-restapi-golang/controllers/v0/feed"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := httprouter.New()

	/*http.HandleFunc("/", index)
	http.HandleFunc("/api/v0/feed", fh.IndexHandler)
	http.ListenAndServe(":"+port, nil)*/

	mux.GET("/", index)
	mux.GET("/api/v0/feed", fh.IndexHandler)
	mux.POST("/api/v0/feed", fh.CreateFeedItemHandler)
	http.ListenAndServe(":"+port, mux)
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "/api/v0/")
}
