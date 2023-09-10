package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func InitilizeDb() {
	// database connection
	var err error
	Db, err = sql.Open("mysql", "root:saurabh123@/htmx_web_server")
	if err != nil {
		log.Fatal("failed to open mysql connection, error: ", err)
	}
	if err := Db.Ping(); err == nil {
		log.Println("ping to db successfull:)")
	}
}

func GetDb() *sql.DB {
	return Db
}
