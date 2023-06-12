package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DB struct {
	conn *sql.DB
}

// CreateDB Returns an instance of a DB connection to be used
func CreateDB(dbHost, dbName, dbUsername, dbPassword string) (*DB, error) {
	database_url := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUsername, dbPassword, dbHost, dbName)
	conn, err := sql.Open("mysql", database_url)
	conn.SetMaxIdleConns(64)
	conn.SetMaxOpenConns(64)
	conn.SetConnMaxLifetime(time.Minute)
	if err != nil {
		return nil, err
	}
	return &DB{conn: conn}, nil
}

// Close closes database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// ExecInsert used for insert query, returns ID of inserted row
func (db *DB) ExecInsert(q string, args ...interface{}) (int64, error) {
	result, err := db.conn.Exec(q, args...)
	if err != nil {
		panic("the END")
		return 0, err
	}
	return result.LastInsertId()
}

// ExecSelectOne used for SELECT query for a single row, returns results
func (db *DB) ExecSelectOne(q string, args ...interface{}) (*sql.Row, error) {
	return db.conn.QueryRow(q, args...), nil
}

// ExecSelectMany used for SELECT query of multiple rows, returns results
func (db *DB) ExecSelectMany(q string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(q, args...)
}

// ExecUpdate used for UPDATE query, returns number of updated rows
func (db *DB) ExecUpdate(q string, args ...interface{}) (int64, error) {
	result, err := db.conn.Exec(q, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
