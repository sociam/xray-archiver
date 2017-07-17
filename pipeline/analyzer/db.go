package main

import (
	"database/sql"
	_ "github.com/lib/pq"

	"fmt"
)

type XrayDb struct {
	db *sql.DB
}

func openDb() (*XrayDb, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%s'",
			cfg.Db.Database, cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port,
		))
	if err != nil {
		return nil, err
	}
	return &XrayDb{db}, nil
}
