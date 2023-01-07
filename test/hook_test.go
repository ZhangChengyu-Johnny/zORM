package test

import (
	"fmt"
	"testing"
	"zORM/log"
	"zORM/session"
)

type Account struct {
	ID       int `zorm:"PRIMIORY KEY`
	Password string
}

func (account *Account) BeforeInsert(s *session.Session) error {
	account.ID += 1000
	log.Info("before insert hook:", account)
	return nil
}

func (account *Account) AfterInsert(s *session.Session) error {
	log.Info("after insert hook:", account)
	return nil
}

func (account *Account) AfterQuery(s *session.Session) error {
	log.Info("after query hook:", account)
	account.Password = "******"
	return nil
}

func TestHook(t *testing.T) {
	s := testInit().Model(&Account{})
	s.DropTable()
	s.CreateTable()

	s.Insert(&Account{1, "123456"}, &Account{2, "654321"}, &Account{3, "13579"})

	u := &Account{}
	s.First(u)

	fmt.Println(u)

}
