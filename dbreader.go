package main

import (
	"database/sql"
	_ "embed"

	_ "github.com/denisenkom/go-mssqldb"
)

//go:embed dbsize.sql
var dbsize string

type DataBase struct {
	DatabaseName string
	TotalSize    float32
	RowSize      float32
	LogSize      float32
	Created      string
	Owner        string
	State        string
	Description  string
}

type DbReader struct {
	db *sql.DB
}

func NewDbReader(connectionString string) (*DbReader, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, err
	}

	return &DbReader{db: db}, nil
}

func (r *DbReader) Close() error {
	return r.db.Close()
}

func (r *DbReader) GetDataBases() (*[]DataBase, error) {
	var databases []DataBase

	rows, err := r.db.Query(dbsize)
	if err != nil {
		return &databases, err
	}
	defer rows.Close()

	for rows.Next() {
		var db DataBase
		if err := rows.Scan(&db.DatabaseName, &db.TotalSize, &db.RowSize, &db.LogSize, &db.Created, &db.Owner, &db.State, &db.Description); err != nil {
			return &databases, err
		}
		databases = append(databases, db)
	}

	return &databases, nil
}
