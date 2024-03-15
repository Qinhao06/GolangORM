package session

import (
	"GolangORM/log"
	"GolangORM/schema"
	"fmt"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	if s.refSchema == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refSchema.Model) {
		s.refSchema = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefSchema() *schema.Schema {
	if s.refSchema == nil {
		log.Error("Session.RefSchema: schema is nil")
	}
	return s.refSchema
}

func (s *Session) CreateTable() error {
	table := s.RefSchema()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ", ")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s)", table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE %s", s.RefSchema().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, args := s.dialect.TableExistSql(s.RefSchema().Name)
	row := s.Raw(sql, args...).QueryRow()
	var tableExist string
	_ = row.Scan(&tableExist)
	return tableExist == s.refSchema.Name
}
