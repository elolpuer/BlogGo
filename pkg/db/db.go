package db

import (
	"database/sql"
	"github.com/elolpuer/Blog/cfg"
	_ "github.com/lib/pq"
)

func ConnectionToDB() (*sql.DB,error) {
	var err error
	db, err := sql.Open("postgres", cfg.GetPostgres())
	if err != nil {
		return nil,err
	}
	if err := db.Ping(); err != nil {
		return nil,err
	}
	return db,nil
}