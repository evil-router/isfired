package database

import (
	"database/sql"
	"fmt"
	"github.com/evil-router/isfired/config"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB = nil

func GetDB() (*sql.DB, error) {
	config := config.Config
	if db == nil {
		conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.DB_User,
			config.DB_Pass,
			config.DB_Host,
			config.DB_Port,
			config.DB_Name)
		d, err := sql.Open("mysql", conn)
		if err != nil {
			return nil, err
		}
		db = d
	}

	return db, nil
}
