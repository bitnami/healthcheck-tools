// Package mysql provides functions for querying mysql databases
package mysql

import (
	"database/sql"
	"fmt"

	// mysql implementation of go's database/sql/driver interface.
	_ "github.com/go-sql-driver/mysql"
)

// Database is a structure that contains the database
// credentials
type Database struct {
	Host string
	Port int
	Name string
	User string
	Pass string
}

// Query is a structure that contains the query database
type Query struct {
	Table  string
	Column string
	Key    string
	Value  string
}

// MySQLQuery returns the result of a MySQL query
func (d Database) MySQLQuery(q Query) (result string, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.User, d.Pass, d.Host, d.Port, d.Name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", err
	}
	defer db.Close()
	if _, err := db.Exec("SET sql_mode='ANSI_QUOTES'"); err != nil {
		return "", err
	}
	query := fmt.Sprintf("SELECT %q FROM %q WHERE %q=?", q.Column, q.Table, q.Key)
	rows, err := db.Query(query, q.Value)
	rows.Next()
	err = rows.Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}
