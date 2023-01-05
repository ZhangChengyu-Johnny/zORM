package test

import (
	"fmt"
	"testing"
	"zORM"
)

func TestD1(t *testing.T) {
	engine, _ := zORM.NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	session := engine.NewSession()
	session.Raw("DROP TABLE IF EXISTS User;").Exec()
	session.Raw("CREATE TABLE User(Name text);").Exec()
	ret, _ := session.Raw("INSERT INTO User(`Name`) VALUES(?), (?)", "Tom", "Sam").Exec()
	count, _ := ret.RowsAffected()
	fmt.Printf("test ok, %d affected\n", count)

}
