package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/badoux/checkmail"
	"github.com/udacity/udagram-restapi-golang/aws"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email         string `json:"email"`
	Password_hash string `json:"password_hash"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type RegisterUserRequest struct {
	Email             string `json:"email"`
	PlainTextPassword string `json:"password"` //eg: imageName.jpeg
}

func GetUserByPk(email string) (User, error) {
	row := aws.DB.QueryRow("SELECT email, password_hash, created_at, updated_at FROM users WHERE email = $1", email)

	item := User{}
	err := row.Scan(&item.Email, &item.Password_hash, &item.CreatedAt, &item.UpdatedAt)
	switch {
	case err == sql.ErrNoRows: //if the error returned is that we did not find anything (ther is no rows)
		return item, errors.New("error finding a user")
	case err != nil: //other types of error, means there is an internal server error (something went wrong with our server)
		log.Printf("Internal server error: %s", err.Error())
		return item, errors.New("error finding a user")
	}

	return item, nil
}

func generatePassword(plainTextPassword string) string {
	bs, _ := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost) //the default cost is 10 rounds
	return string(bs)
}

func RegisterUser(req *http.Request) (interface{}, error) {
	newUser := User{}
	ur := &RegisterUserRequest{}
	err := json.NewDecoder(req.Body).Decode(ur)
	if err != nil {
		return newUser, err
	}

	// check email password valid
	if ur.PlainTextPassword == "" {
		return newUser, errors.New("password is required")
	}

	// check email is valid
	if err := checkmail.ValidateFormat(ur.Email); err != nil {
		return newUser, errors.New("email is required or malformed")
	}

	// find the user
	var temp string
	err = aws.DB.QueryRow("SELECT email FROM users WHERE email = $1", ur.Email).Scan(&temp)
	switch {
	case err == sql.ErrNoRows: //no rows returned meaing the user does not exist
		ph := generatePassword(ur.PlainTextPassword)

		newUser.Email = ur.Email
		newUser.Password_hash = ph
		newUser.CreatedAt = time.Now().Format(time.RFC3339)
		newUser.UpdatedAt = time.Now().Format(time.RFC3339)
	case err != nil: //other types of error, means there is an internal server error (something went wrong with our server)
		log.Printf("Internal server error: %s", err.Error())
		return newUser, err
	default:
		return newUser, errors.New("user may already exist")
	}

	// insert values
	_, err = aws.DB.Exec("INSERT INTO users (email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4)", newUser.Email, newUser.Password_hash, newUser.CreatedAt, newUser.UpdatedAt)
	if err != nil {
		return newUser, errors.New("Internal Server Error." + err.Error())
	}

	return newUser, nil
}
