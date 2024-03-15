package test

import (
	"GolangORM/dialect"
	"GolangORM/engine"
	"GolangORM/log"
	"GolangORM/session"
	"testing"
)

type Account struct {
	ID       int `orm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *session.Session) error {
	log.Info("before inert", account)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *session.Session) error {
	log.Info("after query", account)
	account.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	e, _ := engine.NewEngine("sqlite3", "../test.db", &dialect.Sqlite{})
	s := e.NewSession()
	s.Model(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{1, "123456"}, &Account{2, "qwerty"})

	u := &Account{}

	err := s.First(u)
	if err != nil || u.ID != 1001 || u.Password != "******" {
		t.Fatal("Failed to call hooks after query, got", u)
	}
}
