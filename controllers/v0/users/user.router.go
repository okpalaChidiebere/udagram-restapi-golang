package users

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"] //FYI: the primary key for this table is the email. So we expect an email value as the id
	if !ok {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	u, err := GetUserByPk(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	body, _ := json.Marshal(u) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
