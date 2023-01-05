package test

import (
	"testing"
	"zORM"
)

var (
	user1 = &User{
		Name: "Tom",
		Age:  18,
	}

	user2 = &User{
		Name: "Sam",
		Age:  25,
	}

	user3 = &User{
		Name: "Jack",
		Age:  25,
	}
)

func TestInit(t *testing.T) {
	t.Helper()
	engine, _ := zORM.NewEngine("sqlite3", "test.db")
	s := engine.NewSession().Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test")
	}
}
