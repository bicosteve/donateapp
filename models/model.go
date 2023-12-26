package models

import (
	"database/sql"
	"time"
)

var db *sql.DB

const dbTimeout = time.Second * 3

type Models struct {
	User         User
	JsonResponse JsonResponse
}

func NewConnections(dbPool *sql.DB) Models {
	db = dbPool
	return Models{}

}
