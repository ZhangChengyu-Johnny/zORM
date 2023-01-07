package test

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestDatabasesqlTransaction(t *testing.T) {
	db, _ := sql.Open("sqlite3", "test.db")
	defer db.Close()
	db.Exec("CREATE TABLE IF NOT EXISTS User(`name` text)")

	tx, _ := db.Begin()
	_, err1 := tx.Exec("INSERT INTO User(`Name`) VALUES (?)", "Tom")
	_, err2 := tx.Exec("INSERT INTO User(`Name`) VALUES (?)", "Jack")
	if err1 != nil || err2 != nil {
		fmt.Println("rollback")
		tx.Rollback()
	} else {
		fmt.Println("commit")
		tx.Commit()
	}
}
