package test

import (
	"errors"
	"fmt"
	"testing"
	"zORM"
	"zORM/session"
)

func rollbackTransaction(s *session.Session) (interface{}, error) {
	s.Model(&User{}).CreateTable()
	s.Insert(&User{
		Name: "Tom",
		Age:  18})
	// 手动构造一个错误
	return nil, errors.New("rollback test")
}

func TestRollback(t *testing.T) {
	engine, _ := zORM.NewEngine("sqlite3", "test.db")
	defer engine.Close()
	s := engine.NewSession()
	s.Model(&User{}).DropTable()

	engine.Transaction(rollbackTransaction)
	// 执行事务之后查看表是否存在
	fmt.Println(s.Model(&User{}).HasTable())
}

func commitTransaction(s *session.Session) (interface{}, error) {
	s.Model(&User{}).CreateTable()
	ret, err := s.Insert(&User{Name: "Tom", Age: 666})
	return ret, err
}

func TestCommit(t *testing.T) {
	engine, _ := zORM.NewEngine("sqlite3", "test.db")
	defer engine.Close()
	s := engine.NewSession()
	s.Model(&User{}).DropTable()

	ret, err := engine.Transaction(commitTransaction)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("commit success!", ret)
	}
}
