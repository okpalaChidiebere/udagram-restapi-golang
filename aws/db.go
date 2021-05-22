package aws

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/udacity/udagram-restapi-golang/config"
)

var DB *sql.DB

func init() {
	var err error
	c := config.NewConfig()
	log.Printf("connecting to your database") //the AWS EBS does print logs in the log file you generate from EBS console

	DB, err = sql.Open(c.Dialect, "postgres://"+c.Username+":"+c.Password+"@"+c.Host+"/"+c.Database+"?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}
