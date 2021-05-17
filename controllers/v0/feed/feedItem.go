package feeds

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/udacity/udagram-restapi-golang/aws"
)

type FeedItem struct {
	Id        int64  `json:"id"`
	Caption   string `json:"caption"`
	Url       string `json:"url"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateFeedItemRequest struct {
	Caption string `json:"caption"`
	Url     string `json:"url"`
}

func AllFeedItems() ([]FeedItem, error) {
	rows, err := aws.DB.Query("SELECT id, caption, url, created_at, updated_at FROM feeditem;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fis := make([]FeedItem, 0)
	for rows.Next() {
		fi := FeedItem{}
		err := rows.Scan(&fi.Id, &fi.Caption, &fi.Url, &fi.CreatedAt, &fi.UpdatedAt) // order matters
		if err != nil {
			return nil, err
		}
		fis = append(fis, fi)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return fis, nil
}

func PostFeedItem(req *http.Request) (FeedItem, error) {

	item := FeedItem{}
	ni := &CreateFeedItemRequest{}

	err := json.NewDecoder(req.Body).Decode(ni)
	if err != nil {
		return item, err
	}

	item.Caption = ni.Caption
	item.Url = "https://s3-us-west-1.amazonaws.com/udacity-content/images/icon-eror.svg" //mock value for now
	item.CreatedAt = time.Now().Format(time.RFC3339)
	item.UpdatedAt = time.Now().Format(time.RFC3339)

	// insert values
	err = aws.DB.QueryRow("INSERT INTO feeditem (caption, url, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", item.Caption, item.Url, item.CreatedAt, item.UpdatedAt).Scan(&item.Id)
	if err != nil {
		return item, errors.New("Internal Server Error." + err.Error())
	}
	//learn more about sql transactions with Go here https://stackoverflow.com/questions/40675365/get-back-newly-inserted-row-in-postgres-with-sqlx

	return item, nil
}

func GetFeedItem(id string) (FeedItem, error) {
	row := aws.DB.QueryRow("SELECT id, caption, url, created_at, updated_at FROM feeditem WHERE id = $1", id)

	fi := FeedItem{}
	err := row.Scan(&fi.Id, &fi.Caption, &fi.Url, &fi.CreatedAt, &fi.UpdatedAt)
	switch {
	case err == sql.ErrNoRows: //if the error returned is that we did not find anything (ther is no rows)
		return fi, errors.New("error finding a feed item")
	case err != nil: //other types of error, means there is an internal server error (something went wrong with our server)
		log.Printf("Internal server error: %s", err.Error())
		return fi, errors.New("error finding a feed item")
	}

	return fi, nil
}
