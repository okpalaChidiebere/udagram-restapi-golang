package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Adapter is an alias so I dont have to type so much.
type Adapter func(http.Handler) http.Handler

/*
- the first argument, is the final function to be executed after the middleware functions are ben executed
- the rest of the arguement wil be the middlewared function(s).

You can learn more about the spread(...) operator type in go here
https://yourbasic.org/golang/three-dots-ellipsis/
*/
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	//The adapters/middleware gets executed in the same order as provided in the array.
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

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

func VerificationHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := json.Marshal(map[string]interface{}{
		"auth":    true,
		"message": "Authenticated.",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

/*
This is a middleware function. This method will check if the authorization
header exists and is a valid JWT. If yes, it allows the endpoint to continue
processing. If no, it rejects the request and sends appropriate HTTP status
codes.

We did not implement a complicated authorization system,  but we wanted
to make sure that they are at least authenticated

Can learn more on middleware here
http://dirkwillem.nl/gorilla-mux-middleware/

Remember that http.Handler is an interface!
https://golang.org/pkg/net/http/#Handler
*/
func RequireAuthHandler() Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearerToken := r.Header.Get("Authorization") //the request should have an authorization header

			if len(bearerToken) == 0 { //If there is no header
				JSONError(w, map[string]interface{}{
					"message": "No authorization headers.",
				}, http.StatusUnauthorized) //we return a 401 error
				return //no further function will be executed and the entire flow of middlewares is terminated
			}

			//getting the value of the token from the header
			split := strings.Split(bearerToken, " ")
			if len(split) != 2 {
				JSONError(w, map[string]interface{}{
					"message": "Malformed token.",
				}, http.StatusUnauthorized)
				return
			}

			token := split[1]

			_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				//Make sure that the token method conform to "SigningMethodHMAC"
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(c.Secret), nil
			})
			if err != nil {
				JSONError(w, map[string]interface{}{
					"auth":    false,
					"message": err.Error(),
				}, http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r) //the next method could be another middleware function or the final main function we want to get to
		})
	}
}

//The default, http.Error func returns a plain ext, we had to create our own cstom error to return a JSON
//https://stackoverflow.com/questions/59763852/can-you-return-json-in-golang-http-error
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
