package test

import (
	"GolangORM/dialect"
	"GolangORM/engine"
	"testing"
)

type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {

	e, _ := engine.NewEngine("sqlite3", "../test.db", &dialect.Sqlite{})
	s := e.NewSession()
	s.Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}
