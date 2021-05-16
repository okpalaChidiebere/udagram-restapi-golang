package feeds

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func CreateFeedItemHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	i, err := PostFeedItem(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	body, _ := json.Marshal(i) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
