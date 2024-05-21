package connection

import (
	"core/env"
	"core/util"
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type SQL struct{}

func (ref SQL) SQLite(databaseName string) (*sql.DB, error) {
	var err error

	File := util.File{}

	pwd := env.GetPwd()
	databaseDir := filepath.Join(pwd, "database")
	databaseName = File.AddExtensionIfNotExist(databaseName, "db")
	databasePath := filepath.Join(databaseDir, databaseName)
	err = File.CreateIfNotExist(databasePath)
	if err != nil {
		return nil, fmt.Errorf("error create database file: %v", err)
	}
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("error connect to SQLite: %v", err)
	}
	return db, nil
}

func (ref SQL) MySQL(host string, port int, user string, pass string, name string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, host, port, name))
	if err != nil {
		return nil, fmt.Errorf("error connect to MySQL: %v", err)
	}
	return db, nil
}

func (ref SQL) PostgreSQL(host string, port int, user string, pass string, name string) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name))
	if err != nil {
		return nil, fmt.Errorf("error connect to PostgreSQL: %v", err)
	}
	return db, nil
}

func (ref SQL) MSSQL(host string, port int, user string, pass string, name string) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", host, user, pass, port, name))
	if err != nil {
		return nil, fmt.Errorf("error connect to MSSQL: %v", err)
	}
	return db, nil
}
