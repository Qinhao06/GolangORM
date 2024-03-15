package test

import (
	"GolangORM/dialect"
	"GolangORM/engine"
	"GolangORM/log"
	"GolangORM/session"
	"errors"
	"reflect"
	"testing"
)

func OpenDB(t *testing.T) *engine.Engine {
	t.Helper()

	e, err := engine.NewEngine("sqlite3", "test.db", &dialect.Sqlite{})
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return e
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	e := OpenDB(t)
	defer e.Close()
	s := e.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return nil, errors.New("Error")
	})
	if err.Error() != "Error" || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	e := OpenDB(t)
	defer e.Close()
	s := e.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}

func TestEngine_Migrate(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text PRIMARY KEY, XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	u := &User{}
	for rows.Next() {
		rows.Scan(&u.Name, &u.Age)
		log.Infof("User: %v", u)
	}
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}
