package engine

import (
	"GolangORM/dialect"
	"GolangORM/log"
	"GolangORM/session"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver string, source string, dia dialect.Dialect) (e *Engine, err error) {
	open, err := sql.Open(driver, source)
	if err != nil {
		log.Errorf("open database error: %v", err)
		return
	}
	if err = open.Ping(); err != nil {
		log.Errorf("ping database error: %v", err)
		return
	}
	log.Infof("open database success: %v", source)
	return &Engine{db: open,
		dialect: dia,
	}, nil
}

func (e *Engine) Close() {
	err := e.db.Close()
	if err != nil {
		log.Errorf("close database error: %v", err)
	}
	log.Infof("close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(*session.Session) (interface{}, error)

func (e *Engine) Transaction(fn TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		log.Errorf("begin transaction error: %v", err)
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			if e := s.Rollback(); e != nil {
				log.Errorf("rollback transaction error: %v", e)
			}
			log.Errorf("panic: %v", p)
		} else if err != nil {
			if e := s.Rollback(); e != nil {
				log.Errorf("rollback transaction error: %v", e)
			}
		} else {
			if e := s.Commit(); e != nil {
				log.Errorf("commit transaction error: %v", e)
			}
		}
	}()
	return fn(s)
}

func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if !mapB[v] {
			diff = append(diff, v)
		}
	}
	return
}

func (e *Engine) Migrate(value interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist, create it", s.RefSchema().Name)
			return nil, s.CreateTable()
		}
		table := s.RefSchema()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added columns: %v, deleted columns: %v", addCols, delCols)

		for _, col := range addCols {
			_, err = s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table.Name, col, table.GetField(col).Type)).Exec()
			if err != nil {
				log.Errorf("add column %s error: %v", col, err)
				return
			}
		}

		if len(delCols) == 0 {
			return
		}

		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err

}
