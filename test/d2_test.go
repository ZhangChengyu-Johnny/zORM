package test

import (
	"fmt"
	"testing"
	"zORM"
)

func TestSessionCreateTable(t *testing.T) {
	engine, _ := zORM.NewEngine("sqlite3", "test.db")
	s := engine.NewSession().Model(&User{})
	s.DropTable()
	s.CreateTable()
	fmt.Println(s.HasTable())
}
