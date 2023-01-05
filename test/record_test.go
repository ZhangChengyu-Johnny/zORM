package test

import (
	"fmt"
	"testing"
	"zORM"
	"zORM/session"
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

func testInit() *session.Session {
	engine, _ := zORM.NewEngine("sqlite3", "test.db")
	s := engine.NewSession().Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("failed init test")
	}
	return s
}

func TestFind(t *testing.T) {
	s := testInit()
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}

	for i, u := range users {
		fmt.Println(i, u)
	}

}
