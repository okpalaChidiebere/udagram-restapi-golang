package users

import (
	"encoding/json"
	"net/http"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	i, err := RegisterUser(r)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusNotAcceptable)
		JSONError(w, map[string]interface{}{
			"auth":    false,
			"message": err.Error(),
		}, http.StatusNotAcceptable)
		return
	}

	body, _ := json.Marshal(i) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	i, err := LoginUser(r)
	if err != nil {
		JSONError(w, map[string]interface{}{
			"auth":    false,
			"message": err.Error(),
		}, http.StatusNotAcceptable)
		return
	}

	body, _ := json.Marshal(i) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

//The default, http.Error func returns a plain ext, we had to create our own cstom error to return a JSON
//https://stackoverflow.com/questions/59763852/can-you-return-json-in-golang-http-error
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
