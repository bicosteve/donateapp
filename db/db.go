package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	DB *sql.DB
}

var dbConnection = &DB{}

const maxOpenDbConnections = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

func ConnectMysql(dsn string) (*DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDbConnections)
	db.SetConnMaxIdleTime(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifetime)

	err = testDB(db)

	if err != nil {
		return nil, err
	}

	dbConnection.DB = db

	return dbConnection, nil
}

func testDB(d *sql.DB) error {

	err := d.Ping()

	if err != nil {
		fmt.Printf("Error %v ", err)
		return err
	}

	fmt.Println("DB pinged successfully")

	return nil

}
