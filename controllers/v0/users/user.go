package users

import (
	"database/sql"
	"errors"
	"log"

	"github.com/udacity/udagram-restapi-golang/aws"
)

type User struct {
	Id            int64  `json:"id"`
	Email         string `json:"email"`
	Password_hash string `json:"password_hash"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

func GetUserByPk(id string) (User, error) {
	row := aws.DB.QueryRow("SELECT id, email, password_hash, created_at, updated_at FROM user WHERE id = $1", id)

	item := User{}
	err := row.Scan(&item.Id, &item.Email, &item.Password_hash, &item.CreatedAt, &item.UpdatedAt)
	switch {
	case err == sql.ErrNoRows: //if the error returned is that we did not find anything (ther is no rows)
		return item, errors.New("error finding a user")
	case err != nil: //other types of error, means there is an internal server error (something went wrong with our server)
		log.Printf("Internal server error: %s", err.Error())
		return item, errors.New("error finding a user")
	}

	return item, nil
}
