package feeds

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Get all feed items
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fds, err := AllFeedItems()
	if err != nil {
		log.Printf("Problem getting all feeds: %s", err.Error())

		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(fds) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

//create a feed item
func CreateFeedItemHandler(w http.ResponseWriter, r *http.Request) {
	i, err := PostFeedItem(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	body, _ := json.Marshal(i) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// get a specific feeditem by Primary Key
func GetFeedItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	i, err := GetFeedItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	body, _ := json.Marshal(i) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.Write(body) //we return the record to the client in a sensible payload
}

// Get a signed url to put a new item in the bucket
func GetGetSignedUrlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fn, ok := vars["fileName"]
	if !ok {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	url, err := GetGetSignedUrl(fn)
	if err != nil {
		log.Printf("Problem getting signed url for %s: %s", fn, err.Error())

		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(map[string]interface{}{
		"url": url,
	})
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(body) //we return the record to the client in a sensible payload
}
