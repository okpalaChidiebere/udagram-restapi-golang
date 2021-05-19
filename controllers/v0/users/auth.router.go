package users

import (
	"encoding/json"
	"net/http"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	i, err := RegisterUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	body, _ := json.Marshal(i) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}
