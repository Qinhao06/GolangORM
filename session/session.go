package session

import (
	"GolangORM/clause"
	"GolangORM/dialect"
	"GolangORM/log"
	"GolangORM/schema"
	"database/sql"
	"strings"
)

type Session struct {
	db        *sql.DB
	tx        *sql.Tx
	dialect   dialect.Dialect
	refSchema *schema.Schema
	clause    clause.Clause
	sql       strings.Builder
	sqlArgs   []interface{}
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlArgs = s.sqlArgs[:0]
	s.clause = clause.Clause{}
}

func (s *Session) Raw(sql string, args ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlArgs = append(s.sqlArgs, args...)
	return s
}

func (s *Session) QueryRows() (result *sql.Rows, err error) {
	defer s.Clear()
	log.Infof("SQL: %s, Args: %v", s.sql.String(), s.sqlArgs)
	if result, err = s.DB().Query(s.sql.String(), s.sqlArgs...); err != nil {
		log.Error("Query error: %v", err)
	}
	return
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Infof("SQL: %s, Args: %v", s.sql.String(), s.sqlArgs)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlArgs...); err != nil {
		log.Error("Exec error: %v", err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Infof("SQL: %s, Args: %v", s.sql.String(), s.sqlArgs)
	return s.DB().QueryRow(s.sql.String(), s.sqlArgs...)
}
