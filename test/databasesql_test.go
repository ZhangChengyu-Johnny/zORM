package test

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSQL(t *testing.T) {
	// 连接
	db, _ := sql.Open("sqlite3", "test.db")
	defer func() { db.Close() }()

	// 对表操作
	db.Exec("DROP TABLE IF EXISTS User;")
	db.Exec("CREATE TABLE User(Name text);")

	// 插入操作
	ret, err := db.Exec("INSERT INTO User(`Name`) VALUES (?), (?)", "Tom", "Sam")
	if err == nil {
		affected, _ := ret.RowsAffected()
		log.Println("INSER ID:", affected)
	}

	// 多条记录查询操作
	if rows, err := db.Query("SELECT Name FROM User LIMIT 3;"); err == nil {
		for rows.Next() {
			var name string
			rows.Scan(&name)
			log.Println("SELECT RESULT:", name)
		}
		rows.Close()
	} else {
		log.Println("SELECT FAILED:", err)
	}
}
